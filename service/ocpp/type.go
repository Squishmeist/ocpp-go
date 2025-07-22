package ocpp

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)

type MessageType int
const (
    Request   MessageType = 2
    Confirmation MessageType = 3
)
// IsValid checks if the MessageType is valid.
func (m MessageType) IsValid() bool {
    return m == Request || m == Confirmation
}

type ActionType string
const (
    Heartbeat      ActionType = core.HeartbeatFeatureName
)
// IsValid checks if the ActionType is valid.
func (a ActionType) IsValid() bool {
    return a == Heartbeat
}

// RequestBody represents a request message in the OCPP server.
type RequestBody struct {
    MessageType MessageType                 // 2
    MessageId   string                      // UUID
    Action      ActionType                  // e.g. Heartbeat
    Payload     []byte                      // e.g. interface{}
}
// ConfirmationBody represents a confirmation message in the OCPP server.
type ConfirmationBody struct {
    MessageType MessageType                 // 3
    MessageId   string                      // UUID
    Payload     []byte                      // e.g. interface{}
}

// RequestState represents the state of a request in the OCPP server.
type RequestState struct {
    Type MessageType    // 2
    Id   string         // UUID
    Action ActionType  // e.g. Heartbeat 
}

// State holds the current state of the OCPP server.
type State struct {
    RequestStates []RequestState
}
// Append adds a new RequestState to the State if it doesn't already exist.
func (s *State) Append(state RequestState) {
    for _, existing := range s.RequestStates {
        if existing.Id == state.Id {
            return
        }
    }
    s.RequestStates = append(s.RequestStates, state)
}
// FindId searches for a RequestState by its ID and returns it if found.
func (s *State) FindId(id string) (RequestState, error) {
    for _, existing := range s.RequestStates {
        if existing.Id == id {
            return existing, nil
        }
    }
    return RequestState{}, fmt.Errorf("id %s not found", id)
}
