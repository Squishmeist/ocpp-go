package ocpp

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)

type MessageType string

const (
	Request      MessageType = "request"
	Confirmation MessageType = "confirmation"
)

// checks if the MessageType is valid.
func (m MessageType) IsValid() bool {
	return m == Request || m == Confirmation
}

type ActionType string

const (
	Heartbeat        ActionType = core.HeartbeatFeatureName
	BootNotification ActionType = core.BootNotificationFeatureName
)

// checks if the ActionType is valid.
func (a ActionType) IsValid() bool {
	return a == Heartbeat || a == BootNotification
}

// represents a request message in the OCPP server.
type RequestBody struct {
	Type    MessageType // request
	Uuid    string      // UUID
	Action  ActionType  // e.g. Heartbeat
	Payload []byte      // e.g. interface{}
}

// represents a confirmation message in the OCPP server.
type ConfirmationBody struct {
	Type    MessageType // confirmation
	Uuid    string      // UUID
	Payload []byte      // e.g. interface{}
}

type Pair struct {
	Request      RequestBody
	Confirmation *ConfirmationBody // nil until confirmation arrives
}

// holds the current state of the OCPP server.
type State struct {
	Pairs []Pair
}

// adds a new RequestBody to State, unless a pair with the same MessageId already exists.
func (s *State) AddRequest(req RequestBody) {
	for _, pair := range s.Pairs {
		if pair.Request.Uuid == req.Uuid {
			return // already exists
		}
	}
	s.Pairs = append(s.Pairs, Pair{Request: req})
}

// pairs a ConfirmationBody with its RequestBody by MessageId.
func (s *State) AddConfirmation(conf ConfirmationBody) error {
	for i, pair := range s.Pairs {
		if pair.Request.Uuid == conf.Uuid {
			s.Pairs[i].Confirmation = &conf
			return nil
		}
	}
	return fmt.Errorf("no request found for confirmation id %s", conf.Uuid)
}

// returns the Pair for a given Uuid.
func (s *State) FindByUuid(id string) (*Pair, error) {
	for i, pair := range s.Pairs {
		if pair.Request.Uuid == id {
			return &s.Pairs[i], nil
		}
	}
	return nil, fmt.Errorf("id %s not found", id)
}
