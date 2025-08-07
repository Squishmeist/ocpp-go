package types

import (
	"time"
)

// Represents the tykindpe of message in OCPP.
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

// Represents the action kind in OCPP.
// Defines the actions that can be performed, such as Heartbeat or BootNotification.
type ActionKind string

// Defines the kinds of actions in OCPP.
const (
	Heartbeat        ActionKind = "HEARTBEAT"
	BootNotification ActionKind = "BOOTNOTIFICATION"
)

// Checks if the ActionKind is valid.
func (a ActionKind) IsValid() bool {
	return a == Heartbeat || a == BootNotification
}

// Checks if the ActionKind is valid.
func (a ActionKind) ToPtr() *ActionKind {
	if a.IsValid() {
		return &a
	}
	return nil
}

type Meta struct {
	Id           string
	Serialnumber string
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

type ChargePoint struct {
	ID           string
	Connected    bool
	LastSeen     time.Time
	Availability AvailabilityType
}
