package ocpp

import (
	"context"
	"log/slog"

	"github.com/squishmeist/ocpp-go/internal/core"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func Start(ctx context.Context, topicName, subscriptionName, connectionString string, tp trace.TracerProvider) error {
	client, err := core.NewAzureServiceBusClient(
		core.WithAzureServiceBusServiceName("ocpp"),
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

	machine := NewOcppMachine(
		WithTracerProvider(tp),
	)

	slog.Info("Receiver created successfully", "topic", topicName, "subscription", subscriptionName)

	for {
		messages, err := client.ReceiveMessages(ctx, receiver)
		if err != nil {
			slog.Error("Failed to receive messages", "error", err)
			continue
		}

		for _, msg := range messages {
			slog.Info("Received message", "body", string(msg.Body))

			ctx, span := tp.Tracer("ocpp").Start(ctx, "processMessage", trace.WithAttributes(
				attribute.String("id", msg.MessageID),
				attribute.String("topic", topicName),
				attribute.String("subscription", subscriptionName),
				attribute.String("body", string(msg.Body)),
			))

			err := machine.HandleMessage(ctx, msg.Body)
			if err != nil {
				slog.Error("Error handling message", "error", err)
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				span.End()
				// TODO: Abandon the message in production
				receiver.CompleteMessage(ctx, msg, nil)
				// receiver.AbandonMessage(ctx, msg, nil)
				continue
			}

			slog.Info("state after processing", "state", *machine.state)
			span.SetStatus(codes.Ok, "Message processed successfully")
			span.End()
			// TODO: dont use this in production, this is just for testing
			receiver.CompleteMessage(ctx, msg, nil)
		}
	}
}
