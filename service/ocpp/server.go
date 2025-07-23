package ocpp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/squishmeist/ocpp-go/internal/core"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"
)

func Start(state *State, topicName, subscriptionName, connectionString string, tp *trace.TracerProvider) error {
	tracer := tp.Tracer("ocpp-receiver")
	ctx := context.Background()

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
			slog.Info("Received message", "body", string(msg.Body))

			ctx, span := tracer.Start(ctx, "processMessage")
			span.SetAttributes(
				attribute.String("id", msg.MessageID),
				attribute.String("topic", topicName),
				attribute.String("subscription", subscriptionName),
				attribute.String("body", string(msg.Body)),
			)

			body, err := deconstructBody(ctx, msg.Body)
			if err != nil {
				message := "Failed to deconstruct request body"
				slog.Error(message, "error", err)
				receiver.AbandonMessage(ctx, msg, nil)
				span.RecordError(err)
				span.SetStatus(codes.Error, message)
				span.End()
				continue
			}

			var (
				processErr error
				message    string
			)

			switch b := body.(type) {
			case RequestBody:
				processErr = handleRequestBody(ctx, b, state)
				message = "Error handling RequestBody"
			case ConfirmationBody:
				processErr = handleConfirmationBody(ctx, b, state)
				message = "Error handling ConfirmationBody"
			default:
				processErr = fmt.Errorf("unknown body type")
				message = "Unknown body type"
				slog.Warn(message)
			}

			if processErr != nil {
				receiver.AbandonMessage(ctx, msg, nil)
				slog.Error(message, "error", processErr)
				span.RecordError(processErr)
				span.SetStatus(codes.Error, message)
				span.End()
				continue
			}

			slog.Info("state after processing", "state", *state)
			receiver.CompleteMessage(ctx, msg, nil)
			span.SetStatus(codes.Ok, "Message processed successfully")
			span.End()
		}
	}
}
