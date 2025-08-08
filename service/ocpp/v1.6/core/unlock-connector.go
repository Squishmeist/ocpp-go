package core

import (
	"github.com/go-playground/validator/v10"
	"github.com/squishmeist/ocpp-go/service/ocpp/v1.6/types"
)

// -------------------- Unlock Connector (CS -> CP) --------------------
// Central System can request a Charge Point to unlock a connector. To do so, the Central System SHALL send an UnlockConnectorRequest.
// The purpose of this message: Help EV drivers that have problems unplugging their cable from the Charge Point in case of malfunction of the Connector cable retention.
// When a EV driver calls the CPO help-desk, an operator could manually trigger the sending of an UnlockConnectorRequest to the Charge Point, forcing a new attempt to unlock the connector.
// Hopefully this time the connector unlocks and the EV driver can unplug the cable and drive away.
// The UnlockConnectorRequest SHOULD NOT be used to remotely stop a running transaction, use the Remote Stop Transaction instead.
// Upon receipt of an UnlockConnectorRequest, the Charge Point SHALL respond with a UnlockConnectorConfirmation.
// The response payload SHALL indicate whether the Charge Point was able to unlock its connector.
// If there was a transaction in progress on the specific connector, then Charge Point SHALL finish the transaction first as described in Stop Transaction.
const UnlockConnector = "UnlockConnector"

// Status in response to UnlockConnectorRequest.
type UnlockStatus string

const (
	UnlockStatusUnlocked     UnlockStatus = "Unlocked"
	UnlockStatusUnlockFailed UnlockStatus = "UnlockFailed"
	UnlockStatusNotSupported UnlockStatus = "NotSupported"
)

func isValidUnlockStatus(fl validator.FieldLevel) bool {
	status := UnlockStatus(fl.Field().String())
	switch status {
	case UnlockStatusUnlocked, UnlockStatusUnlockFailed, UnlockStatusNotSupported:
		return true
	default:
		return false
	}
}

// The field definition of the UnlockConnector request payload sent by the Central System to the Charge Point.
type UnlockConnectorRequest struct {
	ConnectorId int `json:"connectorId" validate:"gt=0"`
}

// This field definition of the UnlockConnector confirmation payload, sent by the Charge Point to the Central System in response to an UnlockConnectorRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type UnlockConnectorConfirmation struct {
	Status UnlockStatus `json:"status" validate:"required,unlockStatus16"`
}

func init() {
	_ = types.Validate.RegisterValidation("unlockStatus16", isValidUnlockStatus)
}
