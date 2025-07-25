package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/labstack/echo/v4"
	"github.com/squishmeist/ocpp-go/internal/core"
	"github.com/squishmeist/ocpp-go/internal/core/utils"
	"github.com/squishmeist/ocpp-go/pkg/logging"
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

	server := core.NewHttpServer(
		core.WithHttpServiceName("azure-service-bus"),
	)
	send(server, client, &inbound)

	server.Start(conf.HttpServer.Port)
}

func receive() core.MessageHandler {
	return func(ctx context.Context, topic, subscription string, msg *azservicebus.ReceivedMessage) error {
		slog.Info("Received message", "topic", topic, "subscription", subscription, "body", string(msg.Body))
		return nil
	}
}

func send(server *core.HttpServer, client *core.AzureServiceBusClient, inbound *utils.Topic) {
	server.AddRoute(http.MethodPost, "/send", func(reqCtx echo.Context) error {
		c := reqCtx.Request().Context()
		msgType := reqCtx.QueryParam("msg")

		payload := `[2, "uuid-1", "Heartbeat", {}]` // default

		switch msgType {
		case "requesterror":
			payload = `[2, "uuid-error", "Heartbeat", { "currentTimee": "2025-07-22T11:25:25.230Z" }]`
		case "confirmationerror":
			payload = `[3, "uuid-error", { "currentTimee": "2025-07-22T11:25:25.230Z" }]`
		case "heartbeatrequest":
			payload = `[2, "uuid-1", "Heartbeat", {}]`
		case "heartbeatconfirmation":
			payload = `[3, "uuid-1", { "currentTime": "2025-07-22T11:25:25.230Z" }]`
		case "bootnotificationrequest":
			payload = `[2,"uuid-2", "BootNotification",{
            "chargeBoxSerialNumber": "91234567",
            "chargePointModel": "Zappi",
            "chargePointSerialNumber": "91234567",
            "chargePointVendor": "Myenergi",
            "firmwareVersion": "5540",
            "iccid": "",
            "imsi": "",
            "meterType": "",
            "meterSerialNumber": "91234567"
        }]`
		case "bootnotificationconfirmation":
			payload = `[3,"uuid-2",{
            "currentTime": "2024-04-02T11:44:38Z",
            "interval": 30,
            "status": "Accepted"
        }]`
		}

		if err := client.SendMessage(c, inbound.Name, &azservicebus.Message{
			Body: []byte(payload),
		}); err != nil {
			return reqCtx.String(http.StatusInternalServerError, fmt.Sprintf("Failed to send message: %v", err))
		}

		slog.Info("Message sent successfully", "body", payload)
		return reqCtx.String(http.StatusOK, "Message sent to topic!")
	})

}
