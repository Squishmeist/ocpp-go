package core

import (
	"github.com/go-playground/validator/v10"
	"github.com/squishmeist/ocpp-go/service/ocpp/v1.6/types"
)

// -------------------- Heartbeat (CP -> CS) --------------------
// To let the Central System know that a Charge Point is still connected, a Charge Point sends a heartbeat after a configurable time interval.
// The Charge Point SHALL send a HeartbeatRequest for ensuring that the Central System knows that a Charge Point is still alive.
// Upon receipt of a Heartbeat.req PDU, the Central System SHALL respond with a HeartbeatConfirmation.
// The response payload SHALL contain the current time of the Central System, which is RECOMMENDED to be used by the Charge Point to synchronize its internal clock.
// The Charge Point MAY skip sending a HeartbeatRequest when another payload has been sent to the Central System within the configured heartbeat interval.
// This implies that a Central System SHOULD assume availability of a Charge Point whenever a request has been received, the same way as it would have, when it received a HeartbeatRequest.
// With JSON over WebSocket, sending heartbeats is not mandatory. However, for time synchronization it is advised to at least send one heartbeat per 24 hour.
const Heartbeat = "Heartbeat"

// The field definition of the Heartbeat request payload sent by the Charge Point to the Central System.
type HeartbeatRequest struct {
}

// This field definition of the Heartbeat confirmation payload, sent by the Central System to the Charge Point in response to a HeartbeatRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type HeartbeatConfirmation struct {
	CurrentTime *types.DateTime `json:"currentTime" validate:"required"`
}

func validateHeartbeatConfirmation(sl validator.StructLevel) {
	confirmation := sl.Current().Interface().(HeartbeatConfirmation)
	if types.DateTimeIsNull(confirmation.CurrentTime) {
		sl.ReportError(confirmation.CurrentTime, "CurrentTime", "currentTime", "required", "")
	}
}

func init() {
	types.Validate.RegisterStructValidation(validateHeartbeatConfirmation, HeartbeatConfirmation{})
}
