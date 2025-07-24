package types

import "fmt"

// Holds a RequestBody and an optional ConfirmationBody.
type Pair struct {
	Request      RequestBody
	Confirmation *ConfirmationBody // nil until confirmation arrives
}

// Holds the current state of the OCPP.
type State struct {
	Pairs []Pair
}

// Adds a new RequestBody to state, unless a pair with the same MessageId already exists.
func (s *State) AddRequest(req RequestBody) {
	for _, pair := range s.Pairs {
		if pair.Request.Uuid == req.Uuid {
			return // already exists
		}
	}
	s.Pairs = append(s.Pairs, Pair{Request: req})
}

// Pairs a ConfirmationBody with its RequestBody by MessageId.
func (s *State) AddConfirmation(conf ConfirmationBody) error {
	for i, pair := range s.Pairs {
		if pair.Request.Uuid == conf.Uuid {
			s.Pairs[i].Confirmation = &conf
			return nil
		}
	}
	return fmt.Errorf("no request found for confirmation id %s", conf.Uuid)
}

// Finds the Pair for a given Uuid.
func (s *State) FindByUuid(uuid string) (*Pair, error) {
	for i, pair := range s.Pairs {
		if pair.Request.Uuid == uuid {
			return &s.Pairs[i], nil
		}
	}
	return nil, fmt.Errorf("uuid %s not found", uuid)
}
