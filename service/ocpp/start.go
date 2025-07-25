package ocpp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/squishmeist/ocpp-go/internal/core"
	"github.com/squishmeist/ocpp-go/internal/core/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type OcppStartOption func(*OcppStart)

type OcppStart struct {
	ctx            context.Context
	tracerProvider trace.TracerProvider
	config         utils.Configuration
	client         *core.AzureServiceBusClient
}

func (o *OcppStart) Validate() error {
	if o.tracerProvider == nil {
		return fmt.Errorf("tracer provider is not set")
	}
	if o.ctx == nil {
		return fmt.Errorf("context is not set")
	}
	if o.config == (utils.Configuration{}) {
		return fmt.Errorf("configuration is not set")
	}
	return nil
}

func WithStartContext(ctx context.Context) OcppStartOption {
	return func(o *OcppStart) {
		o.ctx = ctx
	}
}

func WithStartTracerProvider(tp trace.TracerProvider) OcppStartOption {
	return func(o *OcppStart) {
		o.tracerProvider = tp
	}
}

func WithStartConfig(config utils.Configuration) OcppStartOption {
	return func(o *OcppStart) {
		o.config = config
	}
}

func NewOcppStart(opts ...OcppStartOption) *OcppStart {
	start := &OcppStart{}
	for _, opt := range opts {
		opt(start)
	}
	if err := start.Validate(); err != nil {
		slog.Error("Failed to create OcppStart", "error", err)
		panic(err)
	}
	return start
}

// func (o *OcppStart) Start() error {

// }

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

	slog.Info("Sender created successfully", "topic", outbound.Name)

	client.ReceiveMessage(ctx, inbound.Name, inbound.Subscription, getHandler(tp, client, &inbound, &outbound))

	return nil
}

func getHandler(tp trace.TracerProvider, client *core.AzureServiceBusClient, inbound, outbound *utils.Topic) core.MessageHandler {
	machine := NewOcppMachine(
		WithTracerProvider(tp),
	)

	return func(ctx context.Context, topic, subscription string, msg *azservicebus.ReceivedMessage) error {
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
			return err
		}

		if err := client.SendMessage(ctx, outbound.Name, &azservicebus.Message{
			MessageID: &msg.MessageID,
			Body:      []byte(`{"status": "processed", "response": { }}`),
		}); err != nil {
			slog.Error("Failed to send message", "error", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()
			// TODO: dont use this in production
			// receiver.CompleteMessage(ctx, msg, nil)
			return err
		}

		slog.Info("state after processing", "state", *machine.state)
		span.SetStatus(codes.Ok, "Message processed successfully")
		span.End()
		// TODO: dont use this in production
		// receiver.CompleteMessage(ctx, msg, nil)
		return nil
	}
}
