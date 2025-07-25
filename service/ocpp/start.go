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

type OcppOption func(*Ocpp)

type Ocpp struct {
	ctx            context.Context
	tracerProvider trace.TracerProvider
	config         utils.Configuration
	client         *core.AzureServiceBusClient
}

func (o *Ocpp) Validate() error {
	if o.tracerProvider == nil {
		return fmt.Errorf("tracer provider is not set")
	}
	if o.ctx == nil {
		return fmt.Errorf("context is not set")
	}
	if o.config == (utils.Configuration{}) {
		return fmt.Errorf("configuration is not set")
	}
	if o.client == nil {
		return fmt.Errorf("azure service bus client is not set")
	}
	return nil
}

func WithOcppContext(ctx context.Context) OcppOption {
	return func(o *Ocpp) {
		o.ctx = ctx
	}
}

func WithOcppTracerProvider(tp trace.TracerProvider) OcppOption {
	return func(o *Ocpp) {
		o.tracerProvider = tp
	}
}

func WithOcppConfig(config utils.Configuration) OcppOption {
	return func(o *Ocpp) {
		o.config = config
	}
}

func NewOcpp(opts ...OcppOption) *Ocpp {
	start := &Ocpp{}

	for _, opt := range opts {
		opt(start)
	}

	client, err := core.NewAzureServiceBusClient(
		core.WithAzureServiceBusServiceName("ocpp"),
		core.WithAzureServiceBusConnectionString(start.config.AzureServiceBus.ConnectionString),
	)
	if err != nil {
		slog.Error("Failed to create Azure Service Bus client", "error", err)
		panic(err)
	}
	start.client = client

	if err := start.Validate(); err != nil {
		slog.Error("Failed to create Ocpp", "error", err)
		panic(err)
	}

	return start
}

func (o *Ocpp) Start() error {
	defer o.client.Close(o.ctx)
	inbound := o.config.AzureServiceBus.TopicInbound

	handler := o.handler()
	o.client.ReceiveMessage(o.ctx, inbound.Name, inbound.Subscription, handler)

	return nil
}

func (o *Ocpp) handler() core.MessageHandler {
	inbound, outbound := o.config.AzureServiceBus.TopicInbound, o.config.AzureServiceBus.TopicOutbound

	machine := NewOcppMachine(
		WithTracerProvider(o.tracerProvider),
	)

	return func(ctx context.Context, topic, subscription string, msg *azservicebus.ReceivedMessage) error {
		slog.Info("Received message", "body", string(msg.Body))

		ctx, span := o.tracerProvider.Tracer("ocpp").Start(ctx, "processMessage", trace.WithAttributes(
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

		if err := o.client.SendMessage(ctx, outbound.Name, &azservicebus.Message{
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
