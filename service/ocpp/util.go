package ocpp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/squishmeist/ocpp-go/internal/core/util"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// deconstructs a body into a RequestBody.
func deconstructRequestBody(ctx context.Context, uuid string, arr []any) (RequestBody, error) {
	tracer := otel.Tracer("ocpp-receiver")
	_, span := tracer.Start(ctx, "deconstructRequestBody")
	defer span.End()

	errored := func(msg string, err error) (RequestBody, error) {
		return util.ErrWithSpan[RequestBody](span, msg, err)
	}

	if len(arr) < 4 {
		return errored("Invalid request body", fmt.Errorf("expected 4 elements for REQUEST, got %d", len(arr)))
	}

	actionStr, ok := arr[2].(string)
	if !ok {
		return errored("Invalid action type", fmt.Errorf("expected action to be string, got %T", arr[2]))
	}
	action := ActionType(actionStr)
	if !action.IsValid() {
		return errored("Invalid action", fmt.Errorf("invalid action: %s", action))
	}

	payload, err := json.Marshal(arr[3])
	if err != nil {
		return errored("Failed to marshal payload", err)
	}

	span.SetAttributes(
		attribute.String("uuid", uuid),
		attribute.String("type", "request"),
		attribute.String("action", string(action)),
		attribute.String("payload", string(payload)),
	)
	span.SetStatus(codes.Ok, "Successfully deconstructed request body")

	return RequestBody{
		Type:    "request",
		Uuid:    uuid,
		Action:  action,
		Payload: payload,
	}, nil
}

// deconstructs a body into a ConfirmationBody.
func deconstructConfirmationBody(ctx context.Context, uuid string, arr []any) (ConfirmationBody, error) {
	tracer := otel.Tracer("ocpp-receiver")
	_, span := tracer.Start(ctx, "deconstructConfirmationBody")
	defer span.End()

	errored := func(msg string, err error) (ConfirmationBody, error) {
		return util.ErrWithSpan[ConfirmationBody](span, msg, err)
	}

	if len(arr) < 3 {
		return errored("Invalid confirmation body", fmt.Errorf("expected 3 elements for CONFIRMATION, got %d", len(arr)))
	}

	payload, err := json.Marshal(arr[2])
	if err != nil {
		return errored("Failed to marshal payload", err)
	}

	span.SetAttributes(
		attribute.String("uuid", uuid),
		attribute.String("type", "confirmation"),
		attribute.String("payload", string(payload)),
	)
	span.SetStatus(codes.Ok, "Successfully deconstructed confirmation body")

	return ConfirmationBody{
		Type:    "confirmation",
		Uuid:    uuid,
		Payload: payload,
	}, nil
}

// deconstructs a body from a byte slice into a specific structure based on the message type.
func deconstructBody(ctx context.Context, data []byte) (any, error) {
	tracer := otel.Tracer("ocpp-receiver")
	ctx, span := tracer.Start(ctx, "deconstructBody")
	defer span.End()

	errored := func(msg string, err error) (any, error) {
		return util.ErrWithSpan[any](span, msg, err)
	}

	var arr []any
	if err := json.Unmarshal(data, &arr); err != nil {
		return errored("Failed to unmarshal request body", err)
	}
	if len(arr) < 2 {
		return errored("Invalid message format", fmt.Errorf("expected at least 2 elements, got %d", len(arr)))
	}
	msgType, ok := arr[0].(float64)
	if !ok {
		return errored("Invalid message type", fmt.Errorf("expected float64, got %T", arr[0]))
	}
	uuid, ok := arr[1].(string)
	if !ok {
		return errored("Invalid message uuid", fmt.Errorf("expected string, got %T", arr[1]))
	}

	span.SetAttributes(attribute.String("uuid", uuid))

	switch int(msgType) {
	case 2:
		span.SetAttributes(attribute.String("type", "request"))
		result, err := deconstructRequestBody(ctx, uuid, arr)
		if err != nil {
			return errored("Failed to deconstruct request body", err)
		}
		span.SetStatus(codes.Ok, "Successfully deconstructed request body")
		return result, nil
	case 3:
		span.SetAttributes(attribute.String("type", "confirmation"))
		result, err := deconstructConfirmationBody(ctx, uuid, arr)
		if err != nil {
			return errored("Failed to deconstruct confirmation body", err)
		}
		span.SetStatus(codes.Ok, "Successfully deconstructed confirmation body")
		return result, nil
	default:
		return errored("Unknown message type", fmt.Errorf("unknown message type %v", msgType))
	}
}
