package ocpp

import (
	"encoding/json"

	"gopkg.in/go-playground/validator.v9"
)

func asType[T any](data interface{}) (*T, bool) {
    val, ok := data.(*T)
    return val, ok
}

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