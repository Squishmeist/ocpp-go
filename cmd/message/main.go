package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/squishmeist/ocpp-go/internal/core"
	"github.com/squishmeist/ocpp-go/internal/core/utils"
	"github.com/squishmeist/ocpp-go/pkg/logging"
	message "github.com/squishmeist/ocpp-go/service/message"
)

func main() {
	logging.SetupLogger(logging.LevelDebug, logging.LogEnvDevelopment)
	configName := os.Getenv("CONFIG_NAME")
	if configName == "" {
		configName = "azure-service-bus"
	}
	conf := utils.GetConfig("./config", configName, "yaml")
	inbound, outbound := conf.AzureServiceBus.TopicInbound, conf.AzureServiceBus.TopicOutbound

	ctx := context.Background()
	client, err := core.NewAzureServiceBusClient(
		core.WithAzureServiceBusServiceName("azure-service-bus"),
		core.WithAzureServiceBusConnectionString(conf.AzureServiceBus.ConnectionString),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Azure Service Bus client: %v", err))
	}
	defer client.Close(ctx)

	go client.ReceiveMessage(ctx, inbound.Name, inbound.Subscription, receive())
	go client.ReceiveMessage(ctx, outbound.Name, outbound.Subscription, receive())

	server := message.NewServer(conf, client)
	server.Start()
}

func receive() core.MessageHandler {
	return func(ctx context.Context, topic, subscription string, msg *azservicebus.ReceivedMessage) error {
		slog.Info("Received message", "topic", topic, "subscription", subscription, "body", string(msg.Body))
		return nil
	}
}
