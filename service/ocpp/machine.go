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

type OcppStateMachine struct {
	// Store          OcppStore
	// Servicebus     ServiceBus
	TracerProvider trace.TracerProvider
}

func (o *OcppStateMachine) HandleMessage(ctx context.Context, msg []byte) error {
	ctx, span := o.TracerProvider.Tracer("ocpp").Start(ctx, "HandleMessage")
	defer span.End()

	select {
	case <-ctx.Done():
		// TODO: handle context shutdown
		span.RecordError(ctx.Err())
		span.SetStatus(codes.Error, ctx.Err().Error())
		return ctx.Err()
	default:
		span.AddEvent("pre-processing raw", trace.WithAttributes(attribute.String("message", string(msg))))
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
			return o.handleRequest(ctx, *parsedMsg.action, parsedMsg.payload)
		case t.Confirmation:
			return o.handleConfirmation(ctx, parsedMsg.uuid, parsedMsg.payload)
		default:
			return fmt.Errorf("unknown message type")
		}
	}
}

type parsedMessage struct {
	kind    t.MessageKind
	action  *t.ActionKind
	uuid    string
	payload []byte
}

// parses a raw OCPP message into a parsedMessage struct.
func (o *OcppStateMachine) parseRawMessage(msg []byte) (parsedMessage, error) {
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

// parses a message body into a RequestBody.
func (o *OcppStateMachine) parseRequestBody(uuid string, arr []any) (t.RequestBody, error) {
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
		Type:    "request",
		Uuid:    uuid,
		Action:  action,
		Payload: payload,
	}, nil
}

// parses a message body into a ConfirmationBody.
func (o *OcppStateMachine) parseConfirmationBody(uuid string, arr []any) (t.ConfirmationBody, error) {
	if len(arr) < 3 {
		return t.ConfirmationBody{}, fmt.Errorf("invalid confirmation body: expected at least 3 elements for CONFIRMATION, got %d", len(arr))
	}

	payload, err := json.Marshal(arr[2])
	if err != nil {
		return t.ConfirmationBody{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return t.ConfirmationBody{
		Type:    "confirmation",
		Uuid:    uuid,
		Payload: payload,
	}, nil
}

// processes a OCPP request message
func (o *OcppStateMachine) handleRequest(ctx context.Context, action t.ActionKind, payload []byte) error {
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

func (o *OcppStateMachine) handleConfirmation(ctx context.Context, uuid string, payload []byte) error {
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

func (o *OcppStateMachine) getActionKindFromUuid(uuid string) (t.ActionKind, error) {
	match, err := state.FindByUuid(uuid)
	if err != nil {
		return t.ActionKind(""), fmt.Errorf("failed to find request: %w", err)
	}
	return match.Request.Action, nil
}

// OCPP Commands

// TODO: dispatch this heartbeat to somewhere

// TODO: create success or failure response

// TODO: dispatch to service bus on a diff topic - maybe socket-commands

func (o *OcppStateMachine) HandleHeartbeatRequest(ctx context.Context, payload []byte) error {
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

func (o *OcppStateMachine) HandleHeartbeatConfirmation(ctx context.Context, payload []byte) error {
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

func (o *OcppStateMachine) HandleBootNotificationRequest(ctx context.Context, msg any) error {
	return nil
}

func (o *OcppStateMachine) HandleBootNotificationConfirmation(ctx context.Context, msg any) error {
	return nil
}
