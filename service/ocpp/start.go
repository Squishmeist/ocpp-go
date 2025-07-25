package ocpp

import (
	"context"
	"log/slog"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/squishmeist/ocpp-go/internal/core"
	"github.com/squishmeist/ocpp-go/internal/core/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func Start(ctx context.Context, tp trace.TracerProvider, config utils.Configuration) error {
	inbound, outbound := config.AzureServiceBus.TopicInbound, config.AzureServiceBus.TopicOutbound

	client, err := core.NewAzureServiceBusClient(
		core.WithAzureServiceBusServiceName("ocpp"),
		core.WithAzureServiceBusConnectionString(config.AzureServiceBus.ConnectionString),
	)
	if err != nil {
		slog.Error("Failed to create Azure Service Bus client", "error", err)
		panic(err)
	}
	defer client.Close(ctx)

	receiver, err := client.NewReceiverForSubscription(inbound.Name, inbound.Subscription, nil)
	if err != nil {
		slog.Error("Failed to create receiver", "error", err)
		panic(err)
	}
	defer receiver.Close(ctx)
	slog.Info("Receiver created successfully", "topic", inbound.Name)

	sender, err := client.NewSender(outbound.Name, nil)
	if err != nil {
		slog.Error("Failed to create sender", "error", err)
		panic(err)
	}
	defer sender.Close(ctx)
	slog.Info("Sender created successfully", "topic", outbound.Name)

	machine := NewOcppMachine(
		WithTracerProvider(tp),
	)

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
				attribute.String("topic", inbound.Name),
				attribute.String("subscription", inbound.Subscription),
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

			if err := sender.SendMessage(ctx, &azservicebus.Message{
				MessageID: &msg.MessageID,
				Body:      []byte(`{"status": "processed", "response": { }}`),
			}, nil); err != nil {
				slog.Error("Failed to send message", "error", err)
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				span.End()
				// TODO: dont use this in production
				// receiver.CompleteMessage(ctx, msg, nil)
				continue
			}

			slog.Info("state after processing", "state", *machine.state)
			span.SetStatus(codes.Ok, "Message processed successfully")
			span.End()
			// TODO: dont use this in production
			// receiver.CompleteMessage(ctx, msg, nil)
		}
	}
}
