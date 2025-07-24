package messages

import (
	"github.com/go-playground/validator/v10"
	"github.com/squishmeist/ocpp-go/service/ocpp/types"
)

// Type of reset requested by ResetRequest.
type ResetType string

// Result of ResetRequest.
type ResetStatus string

const (
	ResetTypeHard       ResetType   = "Hard"
	ResetTypeSoft       ResetType   = "Soft"
	ResetStatusAccepted ResetStatus = "Accepted"
	ResetStatusRejected ResetStatus = "Rejected"
)

func isValidResetType(fl validator.FieldLevel) bool {
	status := ResetType(fl.Field().String())
	switch status {
	case ResetTypeHard, ResetTypeSoft:
		return true
	default:
		return false
	}
}

func isValidResetStatus(fl validator.FieldLevel) bool {
	status := ResetStatus(fl.Field().String())
	switch status {
	case ResetStatusAccepted, ResetStatusRejected:
		return true
	default:
		return false
	}
}

type ResetRequest struct {
	Type ResetType `json:"type" validate:"required,resetType16"`
}

type ResetConfirmation struct {
	Status ResetStatus `json:"status" validate:"required,resetStatus16"`
}

func init() {
	_ = types.Validate.RegisterValidation("resetType16", isValidResetType)
	_ = types.Validate.RegisterValidation("resetStatus16", isValidResetStatus)
}
