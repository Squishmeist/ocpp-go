package ocpp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/squishmeist/ocpp-go/internal/core"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var state = &State{}

func Start(topicName, subscriptionName, connectionString string, tp trace.TracerProvider) error {
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

			ctx, span := tracer.Start(ctx, "processMessage", trace.WithAttributes(
				attribute.String("id", msg.MessageID),
				attribute.String("topic", topicName),
				attribute.String("subscription", subscriptionName),
				attribute.String("body", string(msg.Body)),
			))

			errored := func(message string, err error) {
				receiver.AbandonMessage(ctx, msg, nil)
				slog.Error(message, "error", err)
				span.RecordError(err)
				span.SetStatus(codes.Error, message)
				span.End()
			}

			body, err := deconstructBody(ctx, msg.Body)
			if err != nil {
				errored("failed to deconstruct body", err)
				continue
			}

			switch b := body.(type) {
			case RequestBody:
				handler, ok := getRequestHandler(b.Action)
				if !ok {
					errored("unknown action", fmt.Errorf("unknown action: %s", b.Action))
					continue
				}
				if err := handler.Handle(RequestHandlerProps{
					ctx:    ctx,
					body:   b,
					tracer: tracer,
				}); err != nil {
					errored(fmt.Sprintf("failed to handle %s request", b.Action), err)
					continue
				}
			case ConfirmationBody:
				handler, ok := getConfirmationHandler(b.Uuid)
				if !ok {
					errored(fmt.Sprintf("unknown uuid: %s", b.Uuid), fmt.Errorf("unknown uuid: %s", b.Uuid))
					continue
				}
				if err := handler.Handle(ConfirmationHandlerProps{
					ctx:    ctx,
					body:   b,
					tracer: tracer,
				}); err != nil {
					errored(fmt.Sprintf("failed to handle %s confirmation", b.Uuid), err)
					continue
				}
			default:
				errored("unknown body type", fmt.Errorf("unknown body type"))
			}

			slog.Info("state after processing", "state", *state)
			receiver.CompleteMessage(ctx, msg, nil)
			span.SetStatus(codes.Ok, "Message processed successfully")
			span.End()
		}
	}
}
