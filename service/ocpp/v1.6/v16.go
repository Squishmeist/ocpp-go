package v16

import (
	"github.com/squishmeist/ocpp-go/service/ocpp/v1.6/core"
	"github.com/squishmeist/ocpp-go/service/ocpp/v1.6/firmware"
	"github.com/squishmeist/ocpp-go/service/ocpp/v1.6/remotetrigger"
)

// Represents the kind of message in OCPP.
// Can be either "request" or "confirmation".
type MessageKind string

// Defines the kind of messages in OCPP.
const (
	Request      MessageKind = "REQUEST"
	Confirmation MessageKind = "CONFIRMATION"
)

// Checks if the MessageKind is valid.
func (m MessageKind) IsValid() bool {
	return m == Request || m == Confirmation
}

type Meta struct {
	Id           string
	Serialnumber string
}

// Represents the action kind in OCPP.
// Defines the actions that can be performed, such as Heartbeat or BootNotification.
type ActionKind string

// Checks if the ActionKind is valid.
func (a ActionKind) IsValid() bool {
	return a == core.Authorize ||
		a == core.BootNotification ||
		a == core.ChangeAvailability ||
		a == core.ChangeConfiguration ||
		a == core.ClearCache ||
		a == firmware.DiagnosticsStatusNotification ||
		a == firmware.FirmwareStatusNotification ||
		a == core.DataTransfer ||
		a == core.GetConfiguration ||
		a == firmware.GetDiagnostics ||
		a == core.Heartbeat ||
		a == core.MeterValues ||
		a == core.RemoteStartTransaction ||
		a == core.RemoteStopTransaction ||
		a == core.Reset ||
		a == core.StartTransaction ||
		a == core.StatusNotification ||
		a == core.StopTransaction ||
		a == remotetrigger.TriggerMessage ||
		a == core.UnlockConnector ||
		a == firmware.UpdateFirmware
}

// Checks if the ActionKind is valid.
func (a ActionKind) ToPtr() *ActionKind {
	if a.IsValid() {
		return &a
	}
	return nil
}

// Represents a Message body in the OCPP.
type MessageBody struct {
	Kind    MessageKind // e.g. REQUEST or CONFIRMATION
	Uuid    string      // UUID
	Action  ActionKind  // e.g. Heartbeat
	Payload []byte
}

// Represents a Request body in the OCPP.
type RequestBody struct {
	Uuid    string     // UUID
	Action  ActionKind // e.g. Heartbeat
	Payload []byte
}

// Represents a confirmation body in the OCPP.
type ConfirmationBody struct {
	Uuid    string // UUID
	Payload []byte // e.g. interface{}
}
