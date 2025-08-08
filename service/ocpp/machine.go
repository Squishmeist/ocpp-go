package ocpp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	v16 "github.com/squishmeist/ocpp-go/service/ocpp/v1.6"
	"github.com/squishmeist/ocpp-go/service/ocpp/v1.6/core"
	"github.com/squishmeist/ocpp-go/service/ocpp/v1.6/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Is a functional option used to configure the OcppMachine.
type OcppMachineOption func(*OcppMachine)

// The main state machine for handling OCPP messages.
type OcppMachine struct {
	// Servicebus     ServiceBus
	TracerProvider trace.TracerProvider
	store          StoreAdapter
	cache          CacheAdapter
}

// Ensures all required fields are set in the OcppMachine.
func (o *OcppMachine) Validate() error {
	if o.TracerProvider == nil {
		return fmt.Errorf("tracer provider is not set")
	}
	if o.store == nil {
		return fmt.Errorf("store is not set")
	}
	if o.cache == nil {
		return fmt.Errorf("cache is not set")
	}
	return nil
}

// Sets the OpenTelemetry tracer provider for the OcppMachine.
func WithTracerProvider(tp trace.TracerProvider) OcppMachineOption {
	return func(m *OcppMachine) {
		m.TracerProvider = tp
	}
}

// Sets the store for the OcppMachine.
func WithStore(store StoreAdapter) OcppMachineOption {
	return func(m *OcppMachine) {
		m.store = store
	}
}

// Sets the cache for the OcppMachine.
func WithCache(cache CacheAdapter) OcppMachineOption {
	return func(m *OcppMachine) {
		m.cache = cache
	}
}

// Creates a new OcppMachine with the provided options.
func NewOcppMachine(opts ...OcppMachineOption) *OcppMachine {
	machine := &OcppMachine{}

	for _, opt := range opts {
		opt(machine)
	}

	if err := machine.Validate(); err != nil {
		slog.Error("Failed to create OcppMachine", "error", err)
		panic(err)
	}

	return machine
}

// Handles an incoming OCPP message.
func (o *OcppMachine) HandleMessage(ctx context.Context, meta v16.Meta, msg []byte) ([]byte, error) {
	ctx, span := o.TracerProvider.Tracer("ocpp").Start(ctx, "HandleMessage")
	defer span.End()

	processed, err := o.cache.HasProcessed(ctx, meta.Id)
	if err != nil {
		err := fmt.Errorf("failed to process message: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if processed {
		slog.Info("Message already processed", "id", meta.Id)
		span.SetStatus(codes.Ok, "Message already processed")
		return nil, nil
	}

	select {
	case <-ctx.Done():
		// TODO: handle context shutdown
		span.RecordError(ctx.Err())
		span.SetStatus(codes.Error, ctx.Err().Error())
		return nil, ctx.Err()
	default:
		parsedMsg, err := o.parseRawMessage(msg)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.AddEvent("error parsing message")
			return nil, err
		}
		span.AddEvent("parsed message")

		switch parsedMsg.kind {
		case v16.Request:
			body, err := o.handleRequest(ctx, meta, parsedMsg)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				span.AddEvent("error marshaling confirmation body")
				return nil, fmt.Errorf("failed to marshal confirmation body: %w", err)
			}
			span.AddEvent("added request message", trace.WithAttributes(
				attribute.String("action", string(*parsedMsg.action)),
				attribute.String("uuid", parsedMsg.uuid),
			))
			return body, nil
		case v16.Confirmation:
			if err := o.handleConfirmation(ctx, meta, parsedMsg.uuid, parsedMsg.payload); err != nil {
				return nil, err
			}
			slog.Info("Paired confirmation message",
				"uuid", parsedMsg.uuid,
				"payload", parsedMsg.payload,
			)
			span.AddEvent("paired confirmation message", trace.WithAttributes(
				attribute.String("uuid", parsedMsg.uuid),
			))
			return nil, nil
		default:
			return nil, fmt.Errorf("unknown message type")
		}
	}
}

// Represents a parsed OCPP message with its kind, action, UUID, and payload.
type parsedMessage struct {
	kind    v16.MessageKind
	action  *v16.ActionKind
	uuid    string
	payload []byte
}

// Parses a raw OCPP message into a parsedMessage struct.
func (o *OcppMachine) parseRawMessage(msg []byte) (parsedMessage, error) {
	var arr []any
	if err := json.Unmarshal(msg, &arr); err != nil {
		return parsedMessage{}, fmt.Errorf("failed to unmarshal request body: %w", err)
	}
	if len(arr) < 2 {
		return parsedMessage{}, fmt.Errorf("invalid message format: expected at least 2 elements, got %d", len(arr))
	}
	msgKind, ok := arr[0].(float64)
	if !ok {
		return parsedMessage{}, fmt.Errorf("invalid message type: expected float64, got %T", arr[0])
	}
	uuid, ok := arr[1].(string)
	if !ok {
		return parsedMessage{}, fmt.Errorf("invalid message uuid: expected string, got %T", arr[1])
	}

	switch int(msgKind) {
	// Request
	case 2:
		result, err := o.parseRequestBody(uuid, arr)
		if err != nil {
			return parsedMessage{}, fmt.Errorf("failed to parse request body: %w", err)
		}
		return parsedMessage{
			kind:    v16.Request,
			action:  &result.Action,
			uuid:    uuid,
			payload: result.Payload,
		}, nil
	// Confirmation
	case 3:
		result, err := o.parseConfirmationBody(uuid, arr)
		if err != nil {
			return parsedMessage{}, fmt.Errorf("failed to parse confirmation body: %w", err)
		}
		return parsedMessage{
			kind:    v16.Confirmation,
			action:  nil,
			uuid:    uuid,
			payload: result.Payload,
		}, nil
	default:
		return parsedMessage{}, fmt.Errorf("unknown message type %f", msgKind)
	}
}

// Parses a message body into a RequestBody.
func (o *OcppMachine) parseRequestBody(uuid string, arr []any) (v16.RequestBody, error) {
	if len(arr) < 4 {
		return v16.RequestBody{}, fmt.Errorf("invalid request body: expected at least 4 elements for REQUEST, got %d", len(arr))
	}

	actionStr, ok := arr[2].(string)
	if !ok {
		return v16.RequestBody{}, fmt.Errorf("invalid action, expected action to be a string, got %T", arr[2])
	}
	action := v16.ActionKind(actionStr)
	if !action.IsValid() {
		return v16.RequestBody{}, fmt.Errorf("invalid action kind: %s", actionStr)
	}

	payload, err := json.Marshal(arr[3])
	if err != nil {
		return v16.RequestBody{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return v16.RequestBody{
		Uuid:    uuid,
		Action:  action,
		Payload: payload,
	}, nil
}

// Parses a message body into a ConfirmationBody.
func (o *OcppMachine) parseConfirmationBody(uuid string, arr []any) (v16.ConfirmationBody, error) {
	if len(arr) < 3 {
		return v16.ConfirmationBody{}, fmt.Errorf("invalid confirmation body: expected at least 3 elements for CONFIRMATION, got %d", len(arr))
	}

	payload, err := json.Marshal(arr[2])
	if err != nil {
		return v16.ConfirmationBody{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return v16.ConfirmationBody{
		Uuid:    uuid,
		Payload: payload,
	}, nil
}

// Processes a OCPP request message
func (o *OcppMachine) handleRequest(ctx context.Context, meta v16.Meta, msg parsedMessage) ([]byte, error) {
	select {
	case <-ctx.Done():
		// TODO: handle context shutdown
		return nil, ctx.Err()
	default:
		var confirmation any
		var err error

		switch *msg.action {
		case v16.ActionKind(core.Heartbeat):
			confirmation, err = o.handleHeartbeatRequest(ctx, meta.Serialnumber, msg.payload)
		case v16.ActionKind(core.BootNotification):
			confirmation, err = o.handleBootNotificationRequest(ctx, msg.payload)
		default:
			return nil, fmt.Errorf("unknown message type")
		}

		if err != nil {
			return nil, err
		}
		body, err := json.Marshal([]any{
			3,
			msg.uuid,
			confirmation,
		})
		slog.Info("Added request message", "body", body)

		if err != nil {
			return nil, err
		}
		return body, nil

	}
}

// Processes a OCPP confirmation message
func (o *OcppMachine) handleConfirmation(ctx context.Context, meta v16.Meta, uuid string, payload []byte) error {
	request, err := o.cache.GetRequestFromUuid(ctx, uuid)
	if err != nil {
		return fmt.Errorf("failed to find request: %v", err)
	}

	switch request.Action {
	default:
		return fmt.Errorf("unknown message type")
	}
}

// OCPP Commands
// TODO: check with router if im allowed to handle this request

// Handles a Heartbeat request.
func (o *OcppMachine) handleHeartbeatRequest(ctx context.Context, serialnumber string, payload []byte) (core.HeartbeatConfirmation, error) {
	var request core.HeartbeatRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		return core.HeartbeatConfirmation{}, err
	}
	if err := types.Validate.Struct(request); err != nil {
		return core.HeartbeatConfirmation{}, err
	}

	confirmation := core.HeartbeatConfirmation{
		CurrentTime: types.Now(),
	}

	if err := o.store.UpdateLastHeartbeat(ctx, serialnumber, confirmation); err != nil {
		return core.HeartbeatConfirmation{}, err
	}

	return confirmation, nil
}

// Handles a BootNotification request.
func (o *OcppMachine) handleBootNotificationRequest(ctx context.Context, payload []byte) (core.BootNotificationConfirmation, error) {
	var request core.BootNotificationRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		return core.BootNotificationConfirmation{}, err
	}
	if err := types.Validate.Struct(request); err != nil {
		return core.BootNotificationConfirmation{}, err
	}

	if err := o.store.AddChargepoint(ctx, request); err != nil {
		return core.BootNotificationConfirmation{}, err
	}

	return core.BootNotificationConfirmation{
		Status:      core.RegistrationStatusAccepted,
		Interval:    30,
		CurrentTime: types.Now(),
	}, nil
}
