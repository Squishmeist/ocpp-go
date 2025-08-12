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

		proxyMode := true // TODO: check with router if allowed to handle message

		switch parsedMsg.kind {
		case v16.Request:
			body, err := o.handleRequest(ctx, proxyMode, meta, parsedMsg)
			if err != nil {
				return nil, err
			}
			span.AddEvent("processed request message", trace.WithAttributes(
				attribute.String("action", string(*parsedMsg.action)),
				attribute.String("uuid", parsedMsg.uuid),
			))
			return body, nil
		case v16.Confirmation:
			if err := o.handleConfirmation(ctx, proxyMode, meta, parsedMsg); err != nil {
				return nil, err
			}
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

// Processes a OCPP request message. If in proxy mode it returns a confirmation to send back else it stores the request in cache to match with a confirmation later.
func (o *OcppMachine) handleRequest(ctx context.Context, proxyMode bool, meta v16.Meta, msg parsedMessage) ([]byte, error) {
	select {
	case <-ctx.Done():
		// TODO: handle context shutdown
		return nil, ctx.Err()
	default:
		var confirmation any
		var err error

		switch *msg.action {
		case v16.ActionKind(core.BootNotification):
			confirmation, err = o.handleBootNotificationRequest(ctx, proxyMode, msg.payload)
		case v16.ActionKind(core.Heartbeat):
			confirmation, err = o.handleHeartbeatRequest(ctx, proxyMode, meta.Serialnumber, msg.payload)
		default:
			return nil, fmt.Errorf("unknown request action")
		}

		if err != nil {
			return nil, err
		}

		// returns confirmation message
		if proxyMode {
			body, err := json.Marshal([]any{
				3,
				msg.uuid,
				confirmation,
			})

			if err != nil {
				return nil, err
			}
			return body, nil
		}

		// stores request in cache
		if err := o.cache.AddRequest(ctx, meta, v16.RequestBody{
			Uuid:    msg.uuid,
			Action:  *msg.action,
			Payload: msg.payload,
		}); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// Processes a OCPP confirmation message. The confirmation is matched with a request in the cache.
func (o *OcppMachine) handleConfirmation(ctx context.Context, proxyMode bool, meta v16.Meta, msg parsedMessage) error {
	request, err := o.cache.GetRequestFromUuid(ctx, msg.uuid)
	if err != nil {
		return err
	}

	switch request.Action {
	case v16.ActionKind(core.BootNotification):
		return o.handleBootNotificationConfirmation(ctx, request, msg.payload)
	case v16.ActionKind(core.Heartbeat):
		return o.handleHeartbeatConfirmation(ctx, meta.Serialnumber, msg.payload)
	default:
		return fmt.Errorf("unknown confirmation action")
	}
}

// Handles a complete BootNotification. AddChargepoint is called to store the Charge Point in the store.
func (o *OcppMachine) onBootNotification(ctx context.Context, request core.BootNotificationRequest) error {
	return o.store.AddChargepoint(ctx, request)
}

// Handles an incoming BootNotification request from a Charge Point.
// Validates the request, processes it via onBootNotification, and returns a confirmation to send if it is in proxy mode.
func (o *OcppMachine) handleBootNotificationRequest(ctx context.Context, proxyMode bool, payload []byte) (core.BootNotificationConfirmation, error) {
	var request core.BootNotificationRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		return core.BootNotificationConfirmation{}, err
	}
	if err := types.Validate.Struct(request); err != nil {
		return core.BootNotificationConfirmation{}, err
	}

	if proxyMode {
		if err := o.onBootNotification(ctx, request); err != nil {
			return core.BootNotificationConfirmation{}, err
		}

		return core.BootNotificationConfirmation{
			Status:      core.RegistrationStatusAccepted,
			Interval:    30,
			CurrentTime: types.Now(),
		}, nil
	}

	return core.BootNotificationConfirmation{}, nil

}

// Handles an incoming BootNotification confirmation from the Central System.
// Validates the confirmation, matches it with the request in cache, and processes via onBootNotification.
func (o *OcppMachine) handleBootNotificationConfirmation(ctx context.Context, request v16.RequestBody, payload []byte) error {
	var confirmation core.BootNotificationConfirmation
	if err := json.Unmarshal(payload, &confirmation); err != nil {
		return err
	}
	if err := types.Validate.Struct(confirmation); err != nil {
		return err
	}

	slog.Debug("Received BootNotification Confirmation",
		slog.Any("payload", confirmation),
	)

	var parsedRequest core.BootNotificationRequest
	if err := json.Unmarshal(request.Payload, &parsedRequest); err != nil {
		return err
	}
	if err := types.Validate.Struct(parsedRequest); err != nil {
		return err
	}

	return o.onBootNotification(ctx, parsedRequest)
}

// Handles a complete Heartbeat. Updates the last heartbeat time in the store.
func (o *OcppMachine) onHeartbeatRequest(ctx context.Context, serialnumber string, confirmation core.HeartbeatConfirmation) error {
	return o.store.UpdateLastHeartbeat(ctx, serialnumber, confirmation)
}

// Handles an incoming Heartbeat request from a Charge Point.
// Validates the request, processes it via onHeartbeatRequest, and returns a confirmation to send if it is in proxy mode.
func (o *OcppMachine) handleHeartbeatRequest(ctx context.Context, proxyMode bool, serialnumber string, payload []byte) (core.HeartbeatConfirmation, error) {
	var request core.HeartbeatRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		return core.HeartbeatConfirmation{}, err
	}
	if err := types.Validate.Struct(request); err != nil {
		return core.HeartbeatConfirmation{}, err
	}

	if proxyMode {
		confirmation := core.HeartbeatConfirmation{
			CurrentTime: types.Now(),
		}

		if err := o.onHeartbeatRequest(ctx, serialnumber, confirmation); err != nil {
			return core.HeartbeatConfirmation{}, err
		}

		return confirmation, nil
	}

	return core.HeartbeatConfirmation{}, nil
}

// Handles an incoming Heartbeat confirmation from the Central System.
// Validates the confirmation, matches it with the request in cache, and processes via onHeartbeatRequest.
func (o *OcppMachine) handleHeartbeatConfirmation(ctx context.Context, serialnumber string, payload []byte) error {
	var confirmation core.HeartbeatConfirmation
	if err := json.Unmarshal(payload, &confirmation); err != nil {
		return err
	}
	if err := types.Validate.Struct(confirmation); err != nil {
		return err
	}

	slog.Debug("Received Heartbeat Confirmation",
		slog.Any("payload", confirmation),
	)

	return o.onHeartbeatRequest(ctx, serialnumber, confirmation)
}
