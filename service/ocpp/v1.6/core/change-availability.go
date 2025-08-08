package core

import (
	"github.com/go-playground/validator/v10"
	"github.com/squishmeist/ocpp-go/service/ocpp/v1.6/types"
)

// -------------------- Change Availability (CS -> CP) --------------------
// Central System can request a Charge Point to change its availability.
// A Charge Point is considered available (“operative”) when it is charging or ready for charging.
// A Charge Point is considered unavailable when it does not allow any charging.
// The Central System SHALL send a ChangeAvailabilityRequest for requesting a Charge Point to change its availability.
// The Central System can change the availability to available or unavailable.
const ChangeAvailability = "ChangeAvailability"

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

// The field definition of the ChangeAvailability request payload sent by the Central System to the Charge Point.
type ChangeAvailabilityRequest struct {
	ConnectorId int              `json:"connectorId" validate:"gte=0"`
	Type        AvailabilityType `json:"type" validate:"required,availabilityType"`
}

// This field definition of the ChangeAvailability confirmation payload, sent by the Charge Point to the Central System in response to a ChangeAvailabilityRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ChangeAvailabilityConfirmation struct {
	Status AvailabilityStatus `json:"status" validate:"required,availabilityStatus"`
}

func init() {
	_ = types.Validate.RegisterValidation("availabilityType", isValidAvailabilityType)
	_ = types.Validate.RegisterValidation("availabilityStatus", isValidAvailabilityStatus)
}
