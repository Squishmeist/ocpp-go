package types

import (
	"github.com/go-playground/validator/v10"
)

// Requested availability change in ChangeAvailabilityRequest.
type AvailabilityType string

const (
	AvailabilityTypeOperative   AvailabilityType = "Operative"
	AvailabilityTypeInoperative AvailabilityType = "Inoperative"
)

func isValidAvailabilityType(fl validator.FieldLevel) bool {
	status := AvailabilityType(fl.Field().String())
	switch status {
	case AvailabilityTypeOperative, AvailabilityTypeInoperative:
		return true
	default:
		return false
	}
}

// Status returned in response to ChangeAvailabilityRequest
type AvailabilityStatus string

const (
	AvailabilityStatusAccepted  AvailabilityStatus = "Accepted"
	AvailabilityStatusRejected  AvailabilityStatus = "Rejected"
	AvailabilityStatusScheduled AvailabilityStatus = "Scheduled"
)

func isValidAvailabilityStatus(fl validator.FieldLevel) bool {
	status := AvailabilityStatus(fl.Field().String())
	switch status {
	case AvailabilityStatusAccepted, AvailabilityStatusRejected, AvailabilityStatusScheduled:
		return true
	default:
		return false
	}
}

type ChangeAvailabilityRequest struct {
	ConnectorId int              `json:"connectorId" validate:"gte=0"`
	Type        AvailabilityType `json:"type" validate:"required,availabilityType"`
}

type ChangeAvailabilityConfirmation struct {
	Status AvailabilityStatus `json:"status" validate:"required,availabilityStatus"`
}

func init() {
	_ = Validate.RegisterValidation("availabilityType", isValidAvailabilityType)
	_ = Validate.RegisterValidation("availabilityStatus", isValidAvailabilityStatus)
}
