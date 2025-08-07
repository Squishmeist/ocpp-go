package message

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/squishmeist/ocpp-go/internal/core"
	ocpppb "github.com/squishmeist/ocpp-go/pkg/api/proto/ocpp/v1"
)

type MessageService struct {
	inboundName string
	client      *core.AzureServiceBusClient
}

type MessageOption func(*MessageService)

func (m *MessageService) Validate() error {
	if m.inboundName == "" {
		return fmt.Errorf("inbound name is not set")
	}
	if m.client == nil {
		return fmt.Errorf("client is not set")
	}

	return nil
}

func WithMessageInboundName(name string) MessageOption {
	return func(m *MessageService) {
		m.inboundName = name
	}
}

func WithMessageClient(client *core.AzureServiceBusClient) MessageOption {
	return func(m *MessageService) {
		m.client = client
	}
}

func NewMessageService(opts ...MessageOption) *MessageService {
	service := &MessageService{}

	for _, opt := range opts {
		opt(service)
	}

	if err := service.Validate(); err != nil {
		slog.Error("failed to validate service", "error", err)
		panic(err)
	}

	return service
}

func (s *MessageService) HeartbeatRequest(ctx context.Context, payload *ocpppb.Request) error {
	return s.client.SendMessage(ctx, s.inboundName, &azservicebus.Message{
		ApplicationProperties: map[string]any{
			"serialnumber": "123456789",
		},
		Body: []byte(`[2, "uuid-1", "Heartbeat", {}]`),
	})
}

func (s *MessageService) HeartbeatConfirmation(ctx context.Context, payload *ocpppb.Request) error {
	return s.client.SendMessage(ctx, s.inboundName, &azservicebus.Message{
		ApplicationProperties: map[string]any{
			"serialnumber": "123456789",
		},
		Body: []byte(`[3, "uuid-1", { "currentTime": "2025-07-22T11:25:25.230Z" }]`),
	})
}

func (s *MessageService) BootNotificationRequest(ctx context.Context, payload *ocpppb.Request) error {
	return s.client.SendMessage(ctx, s.inboundName, &azservicebus.Message{
		ApplicationProperties: map[string]any{
			"serialnumber": "123456789",
		},
		Body: []byte(`[2,"uuid-2", "BootNotification",{
            "chargeBoxSerialNumber": "123456789",
            "chargePointModel": "Zappi",
            "chargePointSerialNumber": "123456789",
            "chargePointVendor": "Myenergi",
            "firmwareVersion": "5540",
            "iccid": "",
            "imsi": "",
            "meterType": "",
            "meterSerialNumber": "91234567"
        }]`),
	})
}

func (s *MessageService) BootNotificationConfirmation(ctx context.Context, payload *ocpppb.Request) error {
	return s.client.SendMessage(ctx, s.inboundName, &azservicebus.Message{
		ApplicationProperties: map[string]any{
			"serialnumber": "123456789",
		},
		Body: []byte(`[3,"uuid-2",{
            "currentTime": "2024-04-02T11:44:38Z",
            "interval": 30,
            "status": "Accepted"
        }]`),
	})
}
