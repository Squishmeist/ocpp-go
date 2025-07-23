package ocpp

import (
	"context"
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
				span.RecordError(err)
				span.SetStatus(codes.Error, message)
				receiver.AbandonMessage(ctx, msg, nil)
				span.End()
				continue
			}

			switch body := body.(type) {
			case RequestBody:
				err := handleRequestBody(ctx, body, state)
				if err != nil {
					message := "Error handling RequestBody"
					slog.Error(message, "error", err)
					span.RecordError(err)
					span.SetStatus(codes.Error, message)
					receiver.AbandonMessage(ctx, msg, nil)
					span.End()
					continue
				}
			case ConfirmationBody:
				err := handleConfirmationBody(ctx, body, state)
				if err != nil {
					message := "Error handling ConfirmationBody"
					slog.Error(message, "error", err)
					span.RecordError(err)
					span.SetStatus(codes.Error, message)
					receiver.AbandonMessage(ctx, msg, nil)
					span.End()
					continue
				}
			default:
				slog.Warn("Unknown body type")
				span.SetStatus(codes.Error, "Unknown body type")
				receiver.AbandonMessage(ctx, msg, nil)
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
