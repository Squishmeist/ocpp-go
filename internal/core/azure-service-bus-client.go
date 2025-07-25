package core

import (
	"context"
	"fmt"

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

func (c *AzureServiceBusClient) NewReceiverForSubscription(topicName, subscriptionName string, opts *azservicebus.ReceiverOptions) (*azservicebus.Receiver, error) {
	return c.Client.NewReceiverForSubscription(topicName, subscriptionName, opts)
}

func (c *AzureServiceBusClient) ReceiveMessages(ctx context.Context, receiver *azservicebus.Receiver) ([]*azservicebus.ReceivedMessage, error) {
	return receiver.ReceiveMessages(ctx, 1, nil)
}

func (c *AzureServiceBusClient) NewSender(queueOrTopic string, opts *azservicebus.NewSenderOptions) (*azservicebus.Sender, error) {
	return c.Client.NewSender(queueOrTopic, opts)
}
