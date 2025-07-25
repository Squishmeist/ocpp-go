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

	sender, err := client.NewSender(inbound.Name, nil)
	if err != nil {
		slog.Error("Failed to create sender", "error", err)
		panic(err)
	}
	defer sender.Close(ctx)
	slog.Info("Sender created successfully", "topic", inbound.Name)

	inboundReceiver, err := client.NewReceiverForSubscription(inbound.Name, inbound.Subscription, nil)
	if err != nil {
		slog.Error("Failed to create receiver", "error", err)
		panic(err)
	}
	defer inboundReceiver.Close(ctx)
	slog.Info("Receiver created successfully", "topic", inbound.Name)

	outboundReceiver, err := client.NewReceiverForSubscription(outbound.Name, outbound.Subscription, nil)
	if err != nil {
		slog.Error("Failed to create receiver", "error", err)
		panic(err)
	}
	defer outboundReceiver.Close(ctx)
	slog.Info("Receiver created successfully", "topic", outbound.Name)

	server := core.NewHttpServer(
		core.WithServiceName("azure-service-bus"),
	)
	send(server, sender)

	go read(ctx, inboundReceiver, "inbound")
	go read(ctx, outboundReceiver, "outbound")

	server.Start(conf.HttpServer.Port)
}

func read(ctx context.Context, receiver *azservicebus.Receiver, direction string) {
	for {
		messages, err := receiver.ReceiveMessages(ctx, 10, nil)
		if err != nil {
			slog.Error("Failed to receive messages", "direction", direction, "error", err)
			continue
		}

		for _, msg := range messages {
			slog.Info("Received message", "direction", direction, "body", string(msg.Body))
		}
	}
}

func send(server *core.HttpServer, sender *azservicebus.Sender) {
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

		if err := sender.SendMessage(c, &azservicebus.Message{
			Body: []byte(payload),
		}, nil); err != nil {
			return reqCtx.String(http.StatusInternalServerError, fmt.Sprintf("Failed to send message: %v", err))
		}

		slog.Info("Message sent successfully", "body", payload)
		return reqCtx.String(http.StatusOK, "Message sent to topic!")
	})

}
