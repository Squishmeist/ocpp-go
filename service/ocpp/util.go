package ocpp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// unmarshals a JSON payload into a struct of type T and validates it.
func unmarshalAndValidate[T any](payload []byte) (*T, error) {
	var obj T
	if err := json.Unmarshal(payload, &obj); err != nil {
		return nil, err
	}
	if err := types.Validate.Struct(obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

// deconstructs a body into a RequestBody.
func deconstructRequestBody(ctx context.Context, uuid string, arr []any) (RequestBody, error) {
	tracer := otel.Tracer("ocpp-receiver")
	_, span := tracer.Start(ctx, "deconstructRequestBody")

	if len(arr) < 4 {
		span.SetStatus(codes.Error, "Invalid request body")
		span.End()
		return RequestBody{}, fmt.Errorf("expected 4 elements for REQUEST, got %d", len(arr))
	}

	actionStr, ok := arr[2].(string)
	if !ok {
		span.SetStatus(codes.Error, "Invalid action type")
		span.End()
		return RequestBody{}, fmt.Errorf("expected action to be string, got %T", arr[2])
	}
	action := ActionType(actionStr)
	if !action.IsValid() {
		span.SetStatus(codes.Error, "Invalid action type")
		span.End()
		return RequestBody{}, fmt.Errorf("invalid action type: %s", action)
	}

	payload, err := json.Marshal(arr[3])
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal payload")
		span.End()
		return RequestBody{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	body := RequestBody{
		Type:    "request",
		Uuid:    uuid,
		Action:  action,
		Payload: payload,
	}

	span.SetAttributes(
		attribute.String("uuid", uuid),
		attribute.String("type", "request"),
		attribute.String("action", string(action)),
	)
	span.End()
	return body, nil
}

// deconstructs a body into a ConfirmationBody.
func deconstructConfirmationBody(uuid string, arr []any) (ConfirmationBody, error) {
	if len(arr) < 3 {
		return ConfirmationBody{}, fmt.Errorf("expected 3 elements for CONFIRMATION, got %d", len(arr))
	}

	payload, err := json.Marshal(arr[2])
	if err != nil {
		return ConfirmationBody{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

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

	var arr []any
	if err := json.Unmarshal(data, &arr); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to unmarshal request body")
		span.End()
		return nil, err
	}
	if len(arr) < 2 {
		span.SetStatus(codes.Error, "Invalid message format")
		span.End()
		return nil, fmt.Errorf("invalid message format")
	}
	msgType, ok := arr[0].(float64)
	if !ok {
		span.SetStatus(codes.Error, "Invalid message type")
		span.End()
		return nil, fmt.Errorf("invalid message type")
	}

	uuid, ok := arr[1].(string)
	if !ok {
		span.SetStatus(codes.Error, "Invalid message uuid")
		span.End()
		return nil, fmt.Errorf("invalid message uuid")
	}

	switch int(msgType) {
	case 2:
		span.SetAttributes(
			attribute.String("uuid", uuid),
			attribute.String("type", "request"),
		)
		result, err := deconstructRequestBody(ctx, uuid, arr)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to deconstruct request body")
			span.End()
			return nil, err
		}
		span.End()
		return result, nil
	case 3:
		return deconstructConfirmationBody(uuid, arr)
	default:
		span.SetStatus(codes.Error, "Unknown message type")
		span.End()
		return nil, fmt.Errorf("unknown message type")
	}
}
