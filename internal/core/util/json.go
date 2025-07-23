package util

import (
	"encoding/json"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

// unmarshals a JSON payload into a struct of type T and validates it.
func UnmarshalAndValidate[T any](payload []byte) (*T, error) {
	var obj T
	if err := json.Unmarshal(payload, &obj); err != nil {
		return nil, err
	}
	if err := types.Validate.Struct(obj); err != nil {
		return nil, err
	}
	return &obj, nil
}
