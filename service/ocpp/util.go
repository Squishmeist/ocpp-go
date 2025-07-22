package ocpp

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

// unmarshals a JSON payload into a struct of type T and validates it.
func unmarshalAndValidate[T any](payload []byte) (*T, error) {
    validate := validator.New()
    var obj T
    if err := json.Unmarshal(payload, &obj); err != nil {
        return nil, err
    }
    if err := validate.Struct(obj); err != nil {
        return nil, err
    }
    return &obj, nil
}

// deconstructs a body into a RequestBody.
func deconstructRequestBody(id string, arr []any) (RequestBody, error) {
    if len(arr) < 4 {
        return RequestBody{}, fmt.Errorf("expected 4 elements for REQUEST, got %d", len(arr))
    }

    actionStr, ok := arr[2].(string)
    if !ok {
        return RequestBody{}, fmt.Errorf("expected action to be string, got %T", arr[2])
    }
    action := ActionType(actionStr)
    if !action.IsValid() {
        return RequestBody{}, fmt.Errorf("invalid action type: %s", action)
    }

    payload, err := json.Marshal(arr[3])
    if err != nil {
        return RequestBody{}, fmt.Errorf("failed to marshal payload: %w", err)
    }

    return RequestBody{
        MessageType: 2,
        MessageId:   id,
        Action:      action,
        Payload:     payload,
    }, nil
}

// deconstructs a body into a ConfirmationBody.
func deconstructConfirmationBody(id string, arr []any) (ConfirmationBody, error) {
    if len(arr) < 3 {
        return ConfirmationBody{}, fmt.Errorf("expected 3 elements for CONFIRMATION, got %d", len(arr))
    }

    payload, err := json.Marshal(arr[2])
    if err != nil {
        return ConfirmationBody{}, fmt.Errorf("failed to marshal payload: %w", err)
    }

    return ConfirmationBody{
        MessageType: 3,
        MessageId:   id,
        Payload:     payload,
    }, nil
}

// deconstructs the context body into a specific structure based on the message type.
func deconstructBody(ctx echo.Context) (any, error) {
	var arr []any
	if err := json.NewDecoder(ctx.Request().Body).Decode(&arr); err != nil {
		return nil, fmt.Errorf("failed to decode request body: %w", err)
	}

	if len(arr) < 2 {
		return nil, fmt.Errorf("expected at least 2 elements, got %d", len(arr))
	}

    bodyType, ok := arr[0].(float64)
    if !ok {
        return nil, fmt.Errorf("expected type to be float64, got %T", arr[0])
    }
    bodyId, ok := arr[1].(string)
    if !ok {
        return nil, fmt.Errorf("expected ID to be string, got %T", arr[1])
    }

    switch int(bodyType) {
    case 2:
        return deconstructRequestBody(bodyId, arr)
    case 3:
        return deconstructConfirmationBody(bodyId, arr)
    default:
        return nil, fmt.Errorf("unsupported message type: %v", bodyType)
    }
}

