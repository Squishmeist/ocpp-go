package types

import (
	"github.com/go-playground/validator/v10"
)

// Reason for stopping a transaction in StopTransactionRequest.
type Reason string

const (
	ReasonDeAuthorized   Reason = "DeAuthorized"
	ReasonEmergencyStop  Reason = "EmergencyStop"
	ReasonEVDisconnected Reason = "EVDisconnected"
	ReasonHardReset      Reason = "HardReset"
	ReasonLocal          Reason = "Local"
	ReasonOther          Reason = "Other"
	ReasonPowerLoss      Reason = "PowerLoss"
	ReasonReboot         Reason = "Reboot"
	ReasonRemote         Reason = "Remote"
	ReasonSoftReset      Reason = "SoftReset"
	ReasonUnlockCommand  Reason = "UnlockCommand"
)

func isValidReason(fl validator.FieldLevel) bool {
	reason := Reason(fl.Field().String())
	switch reason {
	case ReasonDeAuthorized, ReasonEmergencyStop, ReasonEVDisconnected, ReasonHardReset, ReasonLocal, ReasonOther, ReasonPowerLoss, ReasonReboot, ReasonRemote, ReasonSoftReset, ReasonUnlockCommand:
		return true
	default:
		return false
	}
}

type StopTransactionRequest struct {
	IdTag           string       `json:"idTag,omitempty" validate:"max=20"`
	MeterStop       int          `json:"meterStop"`
	Timestamp       *DateTime    `json:"timestamp" validate:"required"`
	TransactionId   int          `json:"transactionId"`
	Reason          Reason       `json:"reason,omitempty" validate:"omitempty,reason"`
	TransactionData []MeterValue `json:"transactionData,omitempty" validate:"omitempty,dive"`
}

type StopTransactionConfirmation struct {
	IdTagInfo *IdTagInfo `json:"idTagInfo,omitempty" validate:"omitempty"`
}

// TODO: advanced validation
func init() {
	_ = Validate.RegisterValidation("reason", isValidReason)
}
