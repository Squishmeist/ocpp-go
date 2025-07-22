package ocpp

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

// asType performs a generic type assertion for interface{} to *T.
func asType[T any](data interface{}) (*T, bool) {
    val, ok := data.(*T)
    return val, ok
}

// unmarshalAndValidate unmarshals a JSON payload into a struct of type T and validates it.
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

// deconstructBody deconstructs the request body into a specific structure based on the message type.
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
        if len(arr) < 4 {
            return nil, fmt.Errorf("expected 4 elements for REQUEST, got %d", len(arr))
        }
        action, ok := arr[2].(string)
        if !ok {
            return nil, fmt.Errorf("expected action to be string, got %T", arr[2])
        }
        payload := arr[3]
        return RequestBody{
            MessageType: 2,
            MessageId:   bodyId,
            Action:      action,
            Payload:     payload,
        }, nil
    case 3:
        if len(arr) < 3 {
            return nil, fmt.Errorf("expected 3 elements for CONFIRMATION, got %d", len(arr))
        }
        payload := arr[2]
        return ConfirmationBody{
            MessageType: 3,
            MessageId:   bodyId,
            Payload:     payload,
        }, nil
    default:
        return nil, fmt.Errorf("unsupported message type: %v", bodyType)
    }
}