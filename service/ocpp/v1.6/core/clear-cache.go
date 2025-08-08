package core

import (
	"github.com/go-playground/validator/v10"
	"github.com/squishmeist/ocpp-go/service/ocpp/v1.6/types"
)

// -------------------- Clear Cache (CS -> CP) --------------------
// Central System can request a Charge Point to clear its Authorization Cache.
// The Central System SHALL send a ClearCacheRequest PDU for clearing the Charge Point’s Authorization Cache.
// Upon receipt of a ClearCacheRequest, the Charge Point SHALL respond with a ClearCacheConfirmation PDU.
// The response PDU SHALL indicate whether the Charge Point was able to clear its Authorization Cache.
const ClearCache = "ClearCache"

// Status returned in response to ClearCacheRequest.
type ClearCacheStatus string

const (
	ClearCacheStatusAccepted ClearCacheStatus = "Accepted"
	ClearCacheStatusRejected ClearCacheStatus = "Rejected"
)

func isValidClearCacheStatus(fl validator.FieldLevel) bool {
	status := ClearCacheStatus(fl.Field().String())
	switch status {
	case ClearCacheStatusAccepted, ClearCacheStatusRejected:
		return true
	default:
		return false
	}
}

// The field definition of the ClearCache request payload sent by the Central System to the Charge Point.
type ClearCacheRequest struct {
}

// This field definition of the ClearCache confirmation payload, sent by the Charge Point to the Central System in response to a ClearCacheRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ClearCacheConfirmation struct {
	Status ClearCacheStatus `json:"status" validate:"required,cacheStatus16"`
}

func init() {
	_ = types.Validate.RegisterValidation("cacheStatus16", isValidClearCacheStatus)
}
