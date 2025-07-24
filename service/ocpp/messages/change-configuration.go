package messages

import (
	"github.com/go-playground/validator/v10"
	"github.com/squishmeist/ocpp-go/service/ocpp/types"
)

const ChangeConfigurationFeatureName = "ChangeConfiguration"

// Status in ChangeConfigurationConfirmation.
type ConfigurationStatus string

const (
	ConfigurationStatusAccepted       ConfigurationStatus = "Accepted"
	ConfigurationStatusRejected       ConfigurationStatus = "Rejected"
	ConfigurationStatusRebootRequired ConfigurationStatus = "RebootRequired"
	ConfigurationStatusNotSupported   ConfigurationStatus = "NotSupported"
)

func isValidConfigurationStatus(fl validator.FieldLevel) bool {
	status := ConfigurationStatus(fl.Field().String())
	switch status {
	case ConfigurationStatusAccepted, ConfigurationStatusRejected, ConfigurationStatusRebootRequired, ConfigurationStatusNotSupported:
		return true
	default:
		return false
	}
}

type ChangeConfigurationRequest struct {
	Key   string `json:"key" validate:"required,max=50"`
	Value string `json:"value" validate:"required,max=500"`
}

type ChangeConfigurationConfirmation struct {
	Status ConfigurationStatus `json:"status" validate:"required,configurationStatus"`
}

func init() {
	_ = types.Validate.RegisterValidation("configurationStatus", isValidConfigurationStatus)
}
