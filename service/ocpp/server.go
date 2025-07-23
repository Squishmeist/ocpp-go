package ocpp

import (
	"context"
	"log/slog"

	"github.com/squishmeist/ocpp-go/internal/core"
	"go.opentelemetry.io/otel/sdk/trace"
)


func Start(state *State, topicName, subscriptionName, connectionString string, tp *trace.TracerProvider) error {
    ctx := context.Background()
    tracer := tp.Tracer("ocpp-receiver")

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
            msgCtx, msgSpan := tracer.Start(ctx, "received-message")

            slog.Info("Received message", "body", string(msg.Body))

            body, err := deconstructBody(msg.Body)
            if err != nil {
                slog.Error("Failed to deconstruct request body:", "error", err)
                receiver.AbandonMessage(msgCtx, msg, nil)
                msgSpan.End()
                continue
            }

            switch body := body.(type) {
            case RequestBody:
                err := handleRequestBody(body, state)
                if err != nil {
                    slog.Error("Error handling RequestBody:", "error", err)
                    receiver.AbandonMessage(msgCtx, msg, nil)
                    msgSpan.AddEvent("Error handling RequestBody")
                    msgSpan.End()
                    continue
                }
            case ConfirmationBody:
                err := handleConfirmationBody(body, state)
                if err != nil {
                    slog.Error("Error handling ConfirmationBody:", "error", err)
                    receiver.AbandonMessage(msgCtx, msg, nil)
                    msgSpan.AddEvent("Error handling ConfirmationBody")
                    msgSpan.End()
                    continue
                }
            default:
                slog.Warn("Unknown body type")
                receiver.AbandonMessage(msgCtx, msg, nil)
                msgSpan.AddEvent("Unknown body type")
                msgSpan.End()
                continue
            }

            slog.Info("state after processing", "state", *state)
            receiver.CompleteMessage(msgCtx, msg, nil)
            msgSpan.AddEvent("Message processed successfully")
            msgSpan.End()
        }
    }
}

