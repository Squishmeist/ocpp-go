package messages

import (
	"github.com/go-playground/validator/v10"
	"github.com/squishmeist/ocpp-go/service/ocpp/types"
)

// Status in response to UnlockConnectorRequest.
type UnlockStatus string

const (
	UnlockStatusUnlocked     UnlockStatus = "Unlocked"
	UnlockStatusUnlockFailed UnlockStatus = "UnlockFailed"
	UnlockStatusNotSupported UnlockStatus = "NotSupported"
)

func isValidUnlockStatus(fl validator.FieldLevel) bool {
	status := UnlockStatus(fl.Field().String())
	switch status {
	case UnlockStatusUnlocked, UnlockStatusUnlockFailed, UnlockStatusNotSupported:
		return true
	default:
		return false
	}
}

type UnlockConnectorRequest struct {
	ConnectorId int `json:"connectorId" validate:"gt=0"`
}

type UnlockConnectorConfirmation struct {
	Status UnlockStatus `json:"status" validate:"required,unlockStatus16"`
}

func init() {
	_ = types.Validate.RegisterValidation("unlockStatus16", isValidUnlockStatus)
}
