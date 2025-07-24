package ocpp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	m "github.com/squishmeist/ocpp-go/service/ocpp/messages"
	t "github.com/squishmeist/ocpp-go/service/ocpp/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Is a functional option used to configure the OcppMachine.
type OcppMachineOption func(*OcppMachine)

// The main state machine for handling OCPP messages.
type OcppMachine struct {
	// Store          OcppStore
	// Servicebus     ServiceBus
	TracerProvider trace.TracerProvider
	state          *t.State
}

// Ensures all required fields are set in the OcppMachine.
func (o *OcppMachine) Validate() error {
	if o.TracerProvider == nil {
		return fmt.Errorf("tracer provider is not set")
	}
	if o.state == nil {
		return fmt.Errorf("state is not set")
	}
	return nil
}

// Sets the OpenTelemetry tracer provider for the OcppMachine.
func WithTracerProvider(tp trace.TracerProvider) OcppMachineOption {
	return func(m *OcppMachine) {
		m.TracerProvider = tp
	}
}

// Sets the state for the OcppMachine.
func WithState(state *t.State) OcppMachineOption {
	return func(m *OcppMachine) {
		m.state = state
	}
}

// Creates a new OcppMachine with the provided options.
func NewOcppMachine(opts ...OcppMachineOption) *OcppMachine {
	machine := &OcppMachine{
		state: &t.State{},
	}

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
func (o *OcppMachine) HandleMessage(ctx context.Context, msg []byte) error {
	ctx, span := o.TracerProvider.Tracer("ocpp").Start(ctx, "HandleMessage")
	defer span.End()

	select {
	case <-ctx.Done():
		// TODO: handle context shutdown
		span.RecordError(ctx.Err())
		span.SetStatus(codes.Error, ctx.Err().Error())
		return ctx.Err()
	default:
		parsedMsg, err := o.parseRawMessage(msg)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.AddEvent("error parsing message")
			return err
		}
		span.AddEvent("parsed message")

		switch parsedMsg.kind {
		case t.Request:
			if err := o.handleRequest(ctx, *parsedMsg.action, parsedMsg.payload); err != nil {
				return err
			}
			o.state.AddRequest(t.RequestBody{
				Action:  *parsedMsg.action,
				Uuid:    parsedMsg.uuid,
				Payload: parsedMsg.payload,
			})
			span.AddEvent("handled request", trace.WithAttributes(
				attribute.String("action", string(*parsedMsg.action)),
				attribute.String("uuid", parsedMsg.uuid),
			))
			return nil
		case t.Confirmation:
			if err := o.handleConfirmation(ctx, parsedMsg.uuid, parsedMsg.payload); err != nil {
				return err
			}
			o.state.AddConfirmation(t.ConfirmationBody{
				Uuid:    parsedMsg.uuid,
				Payload: parsedMsg.payload,
			})
			span.AddEvent("handled confirmation", trace.WithAttributes(
				attribute.String("uuid", parsedMsg.uuid),
			))
			return nil
		default:
			return fmt.Errorf("unknown message type")
		}
	}
}

// Represents a parsed OCPP message with its kind, action, UUID, and payload.
type parsedMessage struct {
	kind    t.MessageKind
	action  *t.ActionKind
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
			kind:    t.Request,
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
			kind:    t.Confirmation,
			action:  nil,
			uuid:    uuid,
			payload: result.Payload,
		}, nil
	default:
		return parsedMessage{}, fmt.Errorf("unknown message type %f", msgKind)
	}
}

// Parses a message body into a RequestBody.
func (o *OcppMachine) parseRequestBody(uuid string, arr []any) (t.RequestBody, error) {
	if len(arr) < 4 {
		return t.RequestBody{}, fmt.Errorf("invalid request body: expected at least 4 elements for REQUEST, got %d", len(arr))
	}

	actionStr, ok := arr[2].(string)
	if !ok {
		return t.RequestBody{}, fmt.Errorf("invalid action, expected action to be a string, got %T", arr[2])
	}
	action := t.ActionKind(actionStr)
	if !action.IsValid() {
		return t.RequestBody{}, fmt.Errorf("invalid action kind: %s", actionStr)
	}

	payload, err := json.Marshal(arr[3])
	if err != nil {
		return t.RequestBody{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return t.RequestBody{
		Uuid:    uuid,
		Action:  action,
		Payload: payload,
	}, nil
}

// Parses a message body into a ConfirmationBody.
func (o *OcppMachine) parseConfirmationBody(uuid string, arr []any) (t.ConfirmationBody, error) {
	if len(arr) < 3 {
		return t.ConfirmationBody{}, fmt.Errorf("invalid confirmation body: expected at least 3 elements for CONFIRMATION, got %d", len(arr))
	}

	payload, err := json.Marshal(arr[2])
	if err != nil {
		return t.ConfirmationBody{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return t.ConfirmationBody{
		Uuid:    uuid,
		Payload: payload,
	}, nil
}

// Processes a OCPP request message
func (o *OcppMachine) handleRequest(ctx context.Context, action t.ActionKind, payload []byte) error {
	select {
	case <-ctx.Done():
		// TODO: handle context shutdown
		return ctx.Err()
	default:
		switch action {
		case t.ActionKind(t.Heartbeat):
			return o.HandleHeartbeatRequest(ctx, payload)
		case t.ActionKind(t.BootNotification):
			return o.HandleBootNotificationRequest(ctx, payload)
		default:
			return fmt.Errorf("unknown message type")
		}
	}
}

// Processes a OCPP confirmation message
func (o *OcppMachine) handleConfirmation(ctx context.Context, uuid string, payload []byte) error {
	actionKind, err := o.getActionKindFromUuid(uuid)
	if err != nil {
		return err
	}

	switch actionKind {
	case t.ActionKind(t.Heartbeat):
		return o.HandleHeartbeatConfirmation(ctx, payload)
	case t.ActionKind(t.BootNotification):
		return o.HandleBootNotificationConfirmation(ctx, payload)
	default:
		return fmt.Errorf("unknown message type")
	}
}

// Retrieves the action kind from the state using the UUID.
func (o *OcppMachine) getActionKindFromUuid(uuid string) (t.ActionKind, error) {
	match, err := o.state.FindByUuid(uuid)
	if err != nil {
		return t.ActionKind(""), fmt.Errorf("failed to find request: %w", err)
	}
	return match.Request.Action, nil
}

// OCPP Commands
// TODO: dispatch this heartbeat to somewhere
// TODO: create success or failure response
// TODO: dispatch to service bus on a diff topic - maybe socket-commands

// Handles a Heartbeat request.
func (o *OcppMachine) HandleHeartbeatRequest(ctx context.Context, payload []byte) error {
	var request m.HeartbeatRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		return err
	}
	if err := t.Validate.Struct(request); err != nil {
		return err
	}

	slog.Debug("Received Heartbeat Request",
		slog.Any("payload", request),
	)
	return nil
}

// Handles a Heartbeat confirmation.
func (o *OcppMachine) HandleHeartbeatConfirmation(ctx context.Context, payload []byte) error {
	var confirmation m.HeartbeatConfirmation
	if err := json.Unmarshal(payload, &confirmation); err != nil {
		return err
	}
	if err := t.Validate.Struct(confirmation); err != nil {
		return err
	}

	slog.Debug("Received Heartbeat Confirmation",
		slog.Any("payload", confirmation),
	)
	return nil
}

// Handles a BootNotification request.
func (o *OcppMachine) HandleBootNotificationRequest(ctx context.Context, payload []byte) error {
	var request m.BootNotificationRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		return err
	}
	if err := t.Validate.Struct(request); err != nil {
		return err
	}

	slog.Debug("Received BootNotification Request",
		slog.Any("payload", request),
	)
	return nil
}

// Handles a BootNotification confirmation.
func (o *OcppMachine) HandleBootNotificationConfirmation(ctx context.Context, payload []byte) error {
	var confirmation m.BootNotificationConfirmation
	if err := json.Unmarshal(payload, &confirmation); err != nil {
		return err
	}
	if err := t.Validate.Struct(confirmation); err != nil {
		return err
	}

	slog.Debug("Received BootNotification Confirmation",
		slog.Any("payload", confirmation),
	)
	return nil
}
