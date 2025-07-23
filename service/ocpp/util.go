package ocpp

import (
	"encoding/json"
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
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

// deconstructs a body from a byte slice into a specific structure based on the message type.
func deconstructBody(data []byte) (any, error) {
    var arr []any
    if err := json.Unmarshal(data, &arr); err != nil {
        return nil, err
    }
    if len(arr) < 2 {
        return nil, fmt.Errorf("invalid message format")
    }
    msgType, ok := arr[0].(float64)
    if !ok {
        return nil, fmt.Errorf("invalid message type")
    }

    id, ok := arr[1].(string)
    if !ok {
        return nil, fmt.Errorf("invalid message id")
    }

    switch int(msgType) {
    case 2:
        return deconstructRequestBody(id, arr)
    case 3:
        return deconstructConfirmationBody(id, arr)
    default:
        return nil, fmt.Errorf("unknown message type")
    }
}

