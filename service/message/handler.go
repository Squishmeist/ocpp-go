package message

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/squishmeist/ocpp-go/internal/core"
	messagepb "github.com/squishmeist/ocpp-go/pkg/api/proto/message/v1"
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

func (s *MessageService) BootNotificationRequest(ctx context.Context, payload *messagepb.Request) error {
	return s.client.SendMessage(ctx, s.inboundName, &azservicebus.Message{
		ApplicationProperties: map[string]any{
			"serialnumber": "123456789",
		},
		Body: []byte(`[2,"uuid-bootNotification", "BootNotification",{
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

func (s *MessageService) BootNotificationConfirmation(ctx context.Context, payload *messagepb.Request) error {
	return s.client.SendMessage(ctx, s.inboundName, &azservicebus.Message{
		ApplicationProperties: map[string]any{
			"serialnumber": "123456789",
		},
		Body: []byte(`[3,"uuid-bootNotification",{
            "currentTime": "2024-04-02T11:44:38Z",
            "interval": 30,
            "status": "Accepted"
        }]`),
	})
}

func (s *MessageService) HeartbeatRequest(ctx context.Context, payload *messagepb.Request) error {
	return s.client.SendMessage(ctx, s.inboundName, &azservicebus.Message{
		ApplicationProperties: map[string]any{
			"serialnumber": "123456789",
		},
		Body: []byte(`[2, "uuid-heartbeat", "Heartbeat", {}]`),
	})
}

func (s *MessageService) HeartbeatConfirmation(ctx context.Context, payload *messagepb.Request) error {
	return s.client.SendMessage(ctx, s.inboundName, &azservicebus.Message{
		ApplicationProperties: map[string]any{
			"serialnumber": "123456789",
		},
		Body: []byte(`[3, "uuid-heartbeat", { "currentTime": "2025-07-22T11:25:25.230Z" }]`),
	})
}

func (s *MessageService) MeterValuesRequest(ctx context.Context, payload *messagepb.Request) error {
	return s.client.SendMessage(ctx, s.inboundName, &azservicebus.Message{
		ApplicationProperties: map[string]any{
			"serialnumber": "123456789",
		},
		Body: []byte(`[2,"uuid-meterValues", "MeterValues",{
			"connectorId": 1,
			"transactionId": 1,
			"meterValue": [{
				"timestamp": "2022-05-20T10:39:19.657Z",
				"sampledValue": 
				[
					{
						"value": "12903.0",
						"context": "Sample.Periodic",
						"format": "Raw",
						"measurand": "Energy.Active.Import.Register",
						"location": "Outlet",
						"unit": "Wh"
					},
					{
						"value": "50000.0",
						"context": "Sample.Periodic",
						"format": "Raw",
						"measurand": "Power.Active.Import",
						"location": "Outlet",
						"unit": "W"
					},
					{
						"value": "50.0",
						"context": "Sample.Periodic",
						"format": "Raw",
						"measurand": "Current.Import",
						"location": "Outlet",
						"unit": "A"
					},
					{
						"value": "63.0",
						"context": "Sample.Periodic",
						"format": "Raw",
						"measurand": "SoC",
						"location": "EV",
						"unit": "Percent"
					},
					{
						"value": "298.8",
						"context": "Sample.Periodic",
						"format": "Raw",
						"measurand": "Temperature",
						"location": "Body",
						"unit": "K"
					},
					{
						"value": "270.4",
						"context": "Sample.Periodic",
						"format": "Raw",
						"measurand": "Voltage",
						"location": "Inlet",
						"unit": "V"
					},
					{
						"value": "60.0",
						"context": "Sample.Periodic",
						"format": "Raw",
						"measurand": "Frequency",
						"location": "Inlet"
					}
				]
			}]
		}]`),
	})
}

func (s *MessageService) MeterValuesConfirmation(ctx context.Context, payload *messagepb.Request) error {
	return s.client.SendMessage(ctx, s.inboundName, &azservicebus.Message{
		ApplicationProperties: map[string]any{
			"serialnumber": "123456789",
		},
		Body: []byte(`[3,"uuid-meterValues",{}]`),
	})
}

func (s *MessageService) StartTransactionRequest(ctx context.Context, payload *messagepb.Request) error {
	return s.client.SendMessage(ctx, s.inboundName, &azservicebus.Message{
		ApplicationProperties: map[string]any{
			"serialnumber": "123456789",
		},
		Body: []byte(`[2,"uuid-startTransaction", "StartTransaction",{
			"connectorId": 1,
			"idTag": "04222182626081",
			"meterStart": 0,
			"timestamp": "2022-06-12T09:13:09.819Z"
        }]`),
	})
}

func (s *MessageService) StartTransactionConfirmation(ctx context.Context, payload *messagepb.Request) error {
	return s.client.SendMessage(ctx, s.inboundName, &azservicebus.Message{
		ApplicationProperties: map[string]any{
			"serialnumber": "123456789",
		},
		Body: []byte(`[3,"uuid-startTransaction",{
			"idTagInfo": {
				"status": "Accepted"
			},
			"transactionId": 1
		}]`),
	})
}

func (s *MessageService) StatusNotificationRequest(ctx context.Context, payload *messagepb.Request) error {
	return s.client.SendMessage(ctx, s.inboundName, &azservicebus.Message{
		ApplicationProperties: map[string]any{
			"serialnumber": "123456789",
		},
		Body: []byte(`[2,"uuid-statusNotification", "StatusNotification",{
			"connectorId": 1,
			"errorCode": "NoError",
			"status": "Preparing",
			"timestamp": "2022-06-12T09:13:00.515Z"
        }]`),
	})
}

func (s *MessageService) StatusNotificationConfirmation(ctx context.Context, payload *messagepb.Request) error {
	return s.client.SendMessage(ctx, s.inboundName, &azservicebus.Message{
		ApplicationProperties: map[string]any{
			"serialnumber": "123456789",
		},
		Body: []byte(`[3,"uuid-statusNotification",{}]`),
	})
}

func (s *MessageService) StopTransactionRequest(ctx context.Context, payload *messagepb.Request) error {
	return s.client.SendMessage(ctx, s.inboundName, &azservicebus.Message{
		ApplicationProperties: map[string]any{
			"serialnumber": "123456789",
		},
		Body: []byte(`[2,"uuid-stopTransaction", "StopTransaction",{
			"reason": "Local",
			"transactionId": 1,
			"meterStop": 4329600,
			"timestamp": "2022-09-08T10:31:26.127Z"
        }]`),
	})
}

func (s *MessageService) StopTransactionConfirmation(ctx context.Context, payload *messagepb.Request) error {
	return s.client.SendMessage(ctx, s.inboundName, &azservicebus.Message{
		ApplicationProperties: map[string]any{
			"serialnumber": "123456789",
		},
		Body: []byte(`[3,"uuid-stopTransaction",{
			"idTagInfo":
			{
				"status": "Accepted"
			}
		}]`),
	})
}
