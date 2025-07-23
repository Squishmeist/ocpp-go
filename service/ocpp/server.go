package ocpp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/squishmeist/ocpp-go/internal/core"
	"go.opentelemetry.io/otel/sdk/trace"
)


func Start(state *State, topicName, subscriptionName, connectionString string, tp *trace.TracerProvider) error {
    ctx := context.Background()
    tracer := tp.Tracer("ocpp-listener")

    client, err := core.NewAzureServiceBusClient(
        core.WithAzureServiceBusServiceName("OCPPService"),
        core.WithAzureServiceBusConnectionString(connectionString),
    )
    if err != nil {
        slog.Error("Failed to create Azure Service Bus client", "error", err)
        panic(err)
    }
    defer client.Close(ctx)

    receiver, err := client.NewReceiverForSubscription(topicName, subscriptionName, nil)
    if err != nil {
        slog.Error("Failed to create receiver", "error", err)
        panic(err)
    }
    defer receiver.Close(ctx)
    slog.Info("Receiver created successfully", "topic", topicName, "subscription", subscriptionName)

    for {
        messages, err := client.ReceiveMessages(ctx, receiver)
        if err != nil {
            slog.Error("Failed to receive messages", "error", err)
            continue
        }

        for _, msg := range messages {
            msgCtx, msgSpan := tracer.Start(ctx, "ProcessMessage")
            fmt.Printf("Received message: %s\n", string(msg.Body))
            // Use your existing deconstructBody logic, but adapt it to accept []byte
            body, err := deconstructBody(msg.Body)
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
            receiver.CompleteMessage(msgCtx, msg, nil)
            msgSpan.End()
        }
    }
}

