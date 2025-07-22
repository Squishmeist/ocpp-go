package ocpp

import (
	"encoding/json"

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