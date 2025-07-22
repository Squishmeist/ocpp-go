package ocpp

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)


func handleRequestBody(body RequestBody, state *State) error {
	switch body.Action {
	case Heartbeat:
		obj, err := unmarshalAndValidate[core.HeartbeatRequest](body.Payload)
		if err != nil {
			return fmt.Errorf("failed to unmarshal HeartbeatRequest: %w", err)
		}
		fmt.Printf("HeartbeatRequest: %v\n", obj)
	default:
		return fmt.Errorf("unknown action: %s", body.Action)
	}

	state.Append(RequestState{Type: body.MessageType, Id: body.MessageId, Action: body.Action})
	return nil
}

func handleConfirmationBody(body ConfirmationBody, state *State) error {
	match, err := state.FindId(body.MessageId)
	if err != nil {
		return fmt.Errorf("RequestState not found: %w", err)
	}

	switch match.Action {
	case Heartbeat:
		obj, err := unmarshalAndValidate[core.HeartbeatConfirmation](body.Payload)
		if err != nil {
			return fmt.Errorf("failed to unmarshal HeartbeatConfirmation: %w", err)
		}
		fmt.Printf("HeartbeatConfirmation: %v\n", obj)
	default:
		return fmt.Errorf("unknown action for confirmation: %s", match.Action)
	}
	return nil
}