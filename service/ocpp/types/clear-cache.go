package types

import (
	"github.com/go-playground/validator/v10"
)

const ClearCacheFeatureName = "ClearCache"

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

type ClearCacheRequest struct {
}

type ClearCacheConfirmation struct {
	Status ClearCacheStatus `json:"status" validate:"required,cacheStatus16"`
}

func init() {
	_ = Validate.RegisterValidation("cacheStatus16", isValidClearCacheStatus)
}
