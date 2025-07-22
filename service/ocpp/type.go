package ocpp

import "fmt"

// RequestBody represents a request message in the OCPP server.
type RequestBody struct {
    MessageType int         // 2
    MessageId   string      // UUID
    Action      string      // e.g. "Heartbeat"
    Payload    	any 		// e.g. interface{}
}

// ConfirmationBody represents a confirmation message in the OCPP server.
type ConfirmationBody struct {
    MessageType int         // 3
    MessageId   string      // UUID
    Payload    	any 		// e.g. interface{}
}

// RequestState represents the state of a request in the OCPP server.
type RequestState struct {
    Type int // 2 or 3
    Id   string      // UUID
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
    return RequestState{}, fmt.Errorf("RequestState with id %s not found", id)
}
