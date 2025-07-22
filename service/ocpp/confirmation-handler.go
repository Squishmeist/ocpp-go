package ocpp

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)

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
		if err := heartbeatConfirmation(body.Payload); err != nil {
			return fmt.Errorf("failed to handle Heartbeat confirmation: %w", err)
		}
	case BootNotification:
		if err := bootnotificationConfirmation(body.Payload); err != nil {
			return fmt.Errorf("failed to handle BootNotification confirmation: %w", err)
		}
	default:
		return fmt.Errorf("unknown action for confirmation: %s", match.Request.Action)
	}

	state.AddConfirmation(ConfirmationBody{
		MessageType: body.MessageType,
		MessageId:   body.MessageId,
		Payload:     body.Payload,
	})

	return nil
}

func heartbeatConfirmation(payload []byte) error {
	obj, err := unmarshalAndValidate[core.HeartbeatConfirmation](payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal HeartbeatConfirmation: %w", err)
	}
	fmt.Printf("HeartbeatConfirmation: %v\n", obj)
	return nil
}

func bootnotificationConfirmation(payload []byte) error {
	obj, err := unmarshalAndValidate[core.BootNotificationConfirmation](payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal BootNotificationConfirmation: %w", err)
	}
	fmt.Printf("BootNotificationConfirmation: %v\n", obj)
	return nil
}