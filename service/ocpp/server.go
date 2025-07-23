package ocpp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/squishmeist/ocpp-go/internal/core"
)


func Start(state *State, topicName, subscriptionName string) error {
    ctx := context.Background()
	connStr := "Endpoint=sb://localhost;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=SAS_KEY_VALUE;UseDevelopmentEmulator=true;"

    client, err := core.NewAzureServiceBusClient(
        core.WithAzureServiceBusServiceName("OCPPService"),
        core.WithAzureServiceBusConnectionString(connStr),
    )
    if err != nil {
        return fmt.Errorf("failed to create service bus client: %w", err)
    }
    defer client.Close(ctx)

    receiver, err := client.NewReceiverForSubscription(topicName, subscriptionName, nil)
    if err != nil {
        return fmt.Errorf("failed to create receiver: %w", err)
    }
    defer receiver.Close(ctx)

    fmt.Printf("Listening to topic '%s' subscription '%s'...\n", topicName, subscriptionName)
    for {
        messages, err := client.ReceiveMessages(ctx, receiver)
        if err != nil {
            fmt.Printf("Error receiving messages: %v\n", err)
            time.Sleep(2 * time.Second)
            continue
        }
        for _, msg := range messages {
            fmt.Printf("Received message: %s\n", string(msg.Body))
            // Use your existing deconstructBody logic, but adapt it to accept []byte
            body, err := deconstructBodyFromBytes(msg.Body)
            if err != nil {
                fmt.Println("Failed to deconstruct request body:", err)
                receiver.AbandonMessage(ctx, msg, nil)
                continue
            }

            switch body := body.(type) {
            case RequestBody:
                err := handleRequestBody(body, state)
                if err != nil {
                    fmt.Println("Error handling RequestBody:", err)
                }
            case ConfirmationBody:
                err := handleConfirmationBody(body, state)
                if err != nil {
                    fmt.Println("Error handling ConfirmationBody:", err)
                }
            default:
                fmt.Println("Unknown body type")
            }

            fmt.Println("State after processing: ", *state)
            receiver.CompleteMessage(context.Background(), msg, nil)
        }
    }
}

func deconstructBodyFromBytes(data []byte) (any, error) {
    // Example: unmarshal as []any, then route to deconstructRequestBody or deconstructConfirmationBody
    var arr []any
    if err := json.Unmarshal(data, &arr); err != nil {
        return nil, err
    }
    if len(arr) < 2 {
        return nil, fmt.Errorf("invalid message format")
    }
    msgType, ok := arr[0].(float64)
    if !ok {
        return nil, fmt.Errorf("invalid message type")
    }
    id, ok := arr[1].(string)
    if !ok {
        return nil, fmt.Errorf("invalid message id")
    }
    switch int(msgType) {
    case 2:
        return deconstructRequestBody(id, arr)
    case 3:
        return deconstructConfirmationBody(id, arr)
    default:
        return nil, fmt.Errorf("unknown message type")
    }
}
