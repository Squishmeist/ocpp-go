package ocpp

import "fmt"


func handleRequestBody(body RequestBody, state *State) {
	fmt.Printf("Received RequestBody: MessageType=%d, MessageId=%s, Action=%s, Payload=%v\n",
		body.MessageType, body.MessageId, body.Action, body.Payload)
	state.Append(RequestState{Type: body.MessageType, Id: body.MessageId})
}

func handleConfirmationBody(body ConfirmationBody, state *State) error {
	fmt.Printf("Received ConfirmationBody: MessageType=%d, MessageId=%s, Payload=%v\n",
		body.MessageType, body.MessageId, body.Payload)
	match, err := state.FindId(body.MessageId)
	if err != nil {
		fmt.Println("Error finding RequestState:", err)
		return fmt.Errorf("RequestState not found: %w", err)
	}
	fmt.Printf("Matched RequestState: Type=%d, Id=%s\n", match.Type, match.Id)
	return nil
}