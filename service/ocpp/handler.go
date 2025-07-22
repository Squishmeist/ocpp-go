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

	state.AddRequest(RequestBody{MessageType: body.MessageType, MessageId: body.MessageId, Action: body.Action})
	return nil
}

func handleConfirmationBody(body ConfirmationBody, state *State) error {
	match, err := state.FindById(body.MessageId)
	if err != nil {
		return fmt.Errorf("RequestBody not found: %w", err)
	}

	if match.Confirmation != nil {
		return fmt.Errorf("confirmation already exists for MessageId: %s", body.MessageId)
	}

	switch match.Request.Action {
	case Heartbeat:
		obj, err := unmarshalAndValidate[core.HeartbeatConfirmation](body.Payload)
		if err != nil {
			return fmt.Errorf("failed to unmarshal HeartbeatConfirmation: %w", err)
		}
		fmt.Printf("HeartbeatConfirmation: %v\n", obj)
		state.AddConfirmation(ConfirmationBody{
			MessageType: body.MessageType,
			MessageId:   body.MessageId,
			Payload:     body.Payload,
		})
	default:
		return fmt.Errorf("unknown action for confirmation: %s", match.Request.Action)
	}
	return nil
}