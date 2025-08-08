package core

import (
	"github.com/go-playground/validator/v10"
	"github.com/squishmeist/ocpp-go/service/ocpp/v1.6/types"
)

// -------------------- Data Transfer (CP -> CS / CS -> CP) --------------------
// If a Charge Point needs to send information to the Central System for a function not supported by OCPP, it SHALL use a DataTransfer message.
// The same functionality may also be offered the other way around, allowing a Central System to send arbitrary custom commands to a Charge Point.
const DataTransfer = "DataTransfer"

// Status in DataTransferConfirmation messages.
type DataTransferStatus string

const (
	DataTransferStatusAccepted         DataTransferStatus = "Accepted"
	DataTransferStatusRejected         DataTransferStatus = "Rejected"
	DataTransferStatusUnknownMessageId DataTransferStatus = "UnknownMessageId"
	DataTransferStatusUnknownVendorId  DataTransferStatus = "UnknownVendorId"
)

func isValidDataTransferStatus(fl validator.FieldLevel) bool {
	status := DataTransferStatus(fl.Field().String())
	switch status {
	case DataTransferStatusAccepted, DataTransferStatusRejected, DataTransferStatusUnknownMessageId, DataTransferStatusUnknownVendorId:
		return true
	default:
		return false
	}
}

// The field definition of the DataTransfer request payload sent by an endpoint to ther other endpoint.
type DataTransferRequest struct {
	VendorId  string      `json:"vendorId" validate:"required,max=255"`
	MessageId string      `json:"messageId,omitempty" validate:"max=50"`
	Data      interface{} `json:"data,omitempty"`
}

// This field definition of the DataTransfer confirmation payload, sent by an endpoint in response to a DataTransferRequest, coming from the other endpoint.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type DataTransferConfirmation struct {
	Status DataTransferStatus `json:"status" validate:"required,dataTransferStatus16"`
	Data   interface{}        `json:"data,omitempty"`
}

func init() {
	_ = types.Validate.RegisterValidation("dataTransferStatus16", isValidDataTransferStatus)
}
