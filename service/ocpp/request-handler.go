package ocpp

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)

func handleRequestBody(body RequestBody, state *State) error {
	switch body.Action {
	case core.HeartbeatFeatureName:
		if err := heartbeatRequest(body.Payload); err != nil {
			return fmt.Errorf("failed to handle Heartbeat request: %w", err)
		}
	case core.BootNotificationFeatureName:
		if err := bootnotificationRequest(body.Payload); err != nil {
			return fmt.Errorf("failed to handle BootNotification request: %w", err)
		}
	default:
		return fmt.Errorf("unknown action: %s", body.Action)
	}

	state.AddRequest(RequestBody{MessageType: body.MessageType, MessageId: body.MessageId, Action: body.Action})
	return nil
}

func heartbeatRequest(payload []byte) error {
	obj, err := unmarshalAndValidate[core.HeartbeatRequest](payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal HeartbeatRequest: %w", err)
	}
	fmt.Printf("HeartbeatRequest: %v\n", obj)
	return nil
}

func bootnotificationRequest(payload []byte) error {
	obj, err := unmarshalAndValidate[core.BootNotificationRequest](payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal BootNotificationRequest: %w", err)
	}
	fmt.Printf("BootNotificationRequest: %v\n", obj)
	return nil
}