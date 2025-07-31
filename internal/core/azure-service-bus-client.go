package core

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type AzureServiceBusClient struct {
	Client           *azservicebus.Client
	ServiceName      string
	connectionString string
}

func (c *AzureServiceBusClient) Validate() error {
	if c.Client == nil {
		return fmt.Errorf("missing required dependency: %s", "Client")
	}
	if c.ServiceName == "" {
		return fmt.Errorf("missing required dependency: %s", "ServiceName")
	}
	if c.connectionString == "" {
		return fmt.Errorf("missing required dependency: %s", "ConnectionString")
	}
	return nil
}

type AzureServiceBusOption func(*AzureServiceBusClient)

func WithAzureServiceBusServiceName(serviceName string) AzureServiceBusOption {
	return func(c *AzureServiceBusClient) {
		c.ServiceName = serviceName
	}
}

func WithAzureServiceBusConnectionString(connectionString string) AzureServiceBusOption {
	return func(c *AzureServiceBusClient) {
		c.connectionString = connectionString
	}
}

func NewAzureServiceBusClient(opts ...AzureServiceBusOption) (*AzureServiceBusClient, error) {
	azureServiceBusClient := &AzureServiceBusClient{}

	for _, opt := range opts {
		opt(azureServiceBusClient)
	}

	client, err := azservicebus.NewClientFromConnectionString(azureServiceBusClient.connectionString, nil)
	if err != nil {
		return nil, err
	}
	azureServiceBusClient.Client = client

	if err := azureServiceBusClient.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate AzureServiceBusClient: %w", err)
	}

	return azureServiceBusClient, nil
}

func (c *AzureServiceBusClient) Close(ctx context.Context) error {
	return c.Client.Close(ctx)
}

func (c *AzureServiceBusClient) SendMessage(ctx context.Context, queueOrTopic string, message *azservicebus.Message) error {
	sender, err := c.Client.NewSender(queueOrTopic, nil)
	if err != nil {
		slog.Error("Failed to create azure service bus sender", "error", err)
		return err
	}

	err = sender.SendMessage(ctx, message, nil)
	if err != nil {
		slog.Error("Failed to send message to topic", "error", err, "queueOrTopic", queueOrTopic)
	}

	return nil
}

type MessageHandler func(ctx context.Context, topic, subscription string, msg *azservicebus.ReceivedMessage) error

func (c *AzureServiceBusClient) ReceiveMessage(
	ctx context.Context,
	topic, subscription string,
	handler MessageHandler,
) error {
	receiver, err := c.Client.NewReceiverForSubscription(topic, subscription, nil)
	if err != nil {
		slog.Error("Failed to create azure service bus receiver", "error", err)
		return err
	}

	for {
		messages, err := receiver.ReceiveMessages(ctx, 10, nil)
		if err != nil {
			slog.Error("Failed to receive messages", "error", err)
			return err
		}

		if len(messages) == 0 {
			slog.Info("No messages received from topic", "topic", topic)
			continue
		}

		slog.Info("Received messages from topic", "topic", topic, "messageCount", len(messages))

		for _, msg := range messages {
			if err := handler(ctx, topic, subscription, msg); err != nil {
				slog.Error("Azure Client, handler failed to handle message", "error", err)
			}
			// TODO: dont use this in production
			receiver.CompleteMessage(ctx, msg, nil)
		}
	}

}
