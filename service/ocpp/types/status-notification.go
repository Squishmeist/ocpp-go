package types

import (
	"github.com/go-playground/validator/v10"
)

// Charge Point status reported in StatusNotificationRequest.
type ChargePointErrorCode string

// Status reported in StatusNotificationRequest.
// A status can be reported for the Charge Point main controller (connectorId = 0) or for a specific connector.
// Status for the Charge Point main controller is a subset of the enumeration: Available, Unavailable or Faulted.
type ChargePointStatus string

const (
	ConnectorLockFailure           ChargePointErrorCode = "ConnectorLockFailure"
	EVCommunicationError           ChargePointErrorCode = "EVCommunicationError"
	GroundFailure                  ChargePointErrorCode = "GroundFailure"
	HighTemperature                ChargePointErrorCode = "HighTemperature"
	InternalError                  ChargePointErrorCode = "InternalError"
	LocalListConflict              ChargePointErrorCode = "LocalListConflict"
	NoError                        ChargePointErrorCode = "NoError"
	OtherError                     ChargePointErrorCode = "OtherError"
	OverCurrentFailure             ChargePointErrorCode = "OverCurrentFailure"
	OverVoltage                    ChargePointErrorCode = "OverVoltage"
	PowerMeterFailure              ChargePointErrorCode = "PowerMeterFailure"
	PowerSwitchFailure             ChargePointErrorCode = "PowerSwitchFailure"
	ReaderFailure                  ChargePointErrorCode = "ReaderFailure"
	ResetFailure                   ChargePointErrorCode = "ResetFailure"
	UnderVoltage                   ChargePointErrorCode = "UnderVoltage"
	WeakSignal                     ChargePointErrorCode = "WeakSignal"
	ChargePointStatusAvailable     ChargePointStatus    = "Available"
	ChargePointStatusPreparing     ChargePointStatus    = "Preparing"
	ChargePointStatusCharging      ChargePointStatus    = "Charging"
	ChargePointStatusSuspendedEVSE ChargePointStatus    = "SuspendedEVSE"
	ChargePointStatusSuspendedEV   ChargePointStatus    = "SuspendedEV"
	ChargePointStatusFinishing     ChargePointStatus    = "Finishing"
	ChargePointStatusReserved      ChargePointStatus    = "Reserved"
	ChargePointStatusUnavailable   ChargePointStatus    = "Unavailable"
	ChargePointStatusFaulted       ChargePointStatus    = "Faulted"
)

func isValidChargePointStatus(fl validator.FieldLevel) bool {
	status := ChargePointStatus(fl.Field().String())
	switch status {
	case ChargePointStatusAvailable, ChargePointStatusPreparing, ChargePointStatusCharging, ChargePointStatusFaulted, ChargePointStatusFinishing, ChargePointStatusReserved, ChargePointStatusSuspendedEV, ChargePointStatusSuspendedEVSE, ChargePointStatusUnavailable:
		return true
	default:
		return false
	}
}

func isValidChargePointErrorCode(fl validator.FieldLevel) bool {
	status := ChargePointErrorCode(fl.Field().String())
	switch status {
	case ConnectorLockFailure, EVCommunicationError, GroundFailure, HighTemperature, InternalError, LocalListConflict, NoError, OtherError, OverVoltage, OverCurrentFailure, PowerMeterFailure, PowerSwitchFailure, ReaderFailure, ResetFailure, UnderVoltage, WeakSignal:
		return true
	default:
		return false
	}
}

type StatusNotificationRequest struct {
	ConnectorId     int                  `json:"connectorId" validate:"gte=0"`
	ErrorCode       ChargePointErrorCode `json:"errorCode" validate:"required,chargePointErrorCode"`
	Info            string               `json:"info,omitempty" validate:"max=50"`
	Status          ChargePointStatus    `json:"status" validate:"required,chargePointStatus"`
	Timestamp       *DateTime            `json:"timestamp,omitempty" validate:"omitempty"`
	VendorId        string               `json:"vendorId,omitempty" validate:"max=255"`
	VendorErrorCode string               `json:"vendorErrorCode,omitempty" validate:"max=50"`
}

type StatusNotificationConfirmation struct {
}

func init() {
	_ = Validate.RegisterValidation("chargePointErrorCode", isValidChargePointErrorCode)
	_ = Validate.RegisterValidation("chargePointStatus", isValidChargePointStatus)
}
