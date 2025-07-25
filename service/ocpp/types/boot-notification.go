package types

import (
	"github.com/go-playground/validator/v10"
)

// Result of registration in response to a BootNotification request.
type RegistrationStatus string

const (
	RegistrationStatusAccepted RegistrationStatus = "Accepted"
	RegistrationStatusPending  RegistrationStatus = "Pending"
	RegistrationStatusRejected RegistrationStatus = "Rejected"
)

func isValidRegistrationStatus(fl validator.FieldLevel) bool {
	status := RegistrationStatus(fl.Field().String())
	switch status {
	case RegistrationStatusAccepted, RegistrationStatusPending, RegistrationStatusRejected:
		return true
	default:
		return false
	}
}

type BootNotificationRequest struct {
	ChargeBoxSerialNumber   string `json:"chargeBoxSerialNumber,omitempty" validate:"max=25"`
	ChargePointModel        string `json:"chargePointModel" validate:"required,max=20"`
	ChargePointSerialNumber string `json:"chargePointSerialNumber,omitempty" validate:"max=25"`
	ChargePointVendor       string `json:"chargePointVendor" validate:"required,max=20"`
	FirmwareVersion         string `json:"firmwareVersion,omitempty" validate:"max=50"`
	Iccid                   string `json:"iccid,omitempty" validate:"max=20"`
	Imsi                    string `json:"imsi,omitempty" validate:"max=20"`
	MeterSerialNumber       string `json:"meterSerialNumber,omitempty" validate:"max=25"`
	MeterType               string `json:"meterType,omitempty" validate:"max=25"`
}

type BootNotificationConfirmation struct {
	CurrentTime *DateTime          `json:"currentTime" validate:"required"`
	Interval    int                `json:"interval" validate:"gte=0"`
	Status      RegistrationStatus `json:"status" validate:"required,registrationStatus16"`
}

func (h *BootNotificationRequest) Validate() error {
	if err := Validate.Struct(h); err != nil {
		return err
	}
	return nil
}

func init() {
	_ = Validate.RegisterValidation("registrationStatus16", isValidRegistrationStatus)
}
