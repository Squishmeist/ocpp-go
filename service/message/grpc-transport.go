package message

import (
	"context"

	messagepb "github.com/squishmeist/ocpp-go/pkg/api/proto/message/v1"
	"go.opentelemetry.io/otel"
)

type MessageServiceInterface interface {
	BootNotificationRequest(context.Context, *messagepb.Request) error
	BootNotificationConfirmation(context.Context, *messagepb.Request) error
	HeartbeatRequest(context.Context, *messagepb.Request) error
	HeartbeatConfirmation(context.Context, *messagepb.Request) error
	MeterValuesRequest(context.Context, *messagepb.Request) error
	MeterValuesConfirmation(context.Context, *messagepb.Request) error
	StartTransactionRequest(context.Context, *messagepb.Request) error
	StartTransactionConfirmation(context.Context, *messagepb.Request) error
	StatusNotificationRequest(context.Context, *messagepb.Request) error
	StatusNotificationConfirmation(context.Context, *messagepb.Request) error
	StopTransactionRequest(context.Context, *messagepb.Request) error
	StopTransactionConfirmation(context.Context, *messagepb.Request) error
}

type MessageGrpcTransport struct {
	messagepb.UnimplementedOCPPMessageServer
	Service MessageServiceInterface
}

func NewMessageGrpcTransport(service MessageServiceInterface) *MessageGrpcTransport {
	return &MessageGrpcTransport{
		Service: service,
	}
}

func (h *MessageGrpcTransport) BootNotificationRequest(ctx context.Context, req *messagepb.Request) (*messagepb.Response, error) {
	tracer := otel.Tracer("ocpp-go/service/message")
	ctx, span := tracer.Start(ctx, "BootNotificationRequest")
	defer span.End()

	if err := h.Service.BootNotificationRequest(ctx, req); err != nil {
		return &messagepb.Response{
			Message: "Failed to send boot notification request",
		}, err
	}

	return &messagepb.Response{
		Message: "Boot notification request sent",
	}, nil
}

func (h *MessageGrpcTransport) BootNotificationConfirmation(ctx context.Context, req *messagepb.Request) (*messagepb.Response, error) {
	tracer := otel.Tracer("ocpp-go/service/message")
	ctx, span := tracer.Start(ctx, "BootNotificationConfirmation")
	defer span.End()

	if err := h.Service.BootNotificationConfirmation(ctx, req); err != nil {
		return &messagepb.Response{
			Message: "Failed to send boot notification confirmation",
		}, err
	}

	return &messagepb.Response{
		Message: "Boot notification confirmation sent",
	}, nil
}

func (h *MessageGrpcTransport) HeartbeatRequest(ctx context.Context, req *messagepb.Request) (*messagepb.Response, error) {
	tracer := otel.Tracer("ocpp-go/service/message")
	ctx, span := tracer.Start(ctx, "HeartbeatRequest")
	defer span.End()

	if err := h.Service.HeartbeatRequest(ctx, req); err != nil {
		return &messagepb.Response{
			Message: "Failed to send heartbeat request",
		}, err
	}

	return &messagepb.Response{
		Message: "Heartbeat request sent",
	}, nil
}

func (h *MessageGrpcTransport) HeartbeatConfirmation(ctx context.Context, req *messagepb.Request) (*messagepb.Response, error) {
	tracer := otel.Tracer("ocpp-go/service/message")
	ctx, span := tracer.Start(ctx, "HeartbeatConfirmation")
	defer span.End()

	if err := h.Service.HeartbeatConfirmation(ctx, req); err != nil {
		return &messagepb.Response{
			Message: "Failed to send heartbeat confirmation",
		}, err
	}

	return &messagepb.Response{
		Message: "Heartbeat confirmation sent",
	}, nil
}

func (h *MessageGrpcTransport) MeterValuesRequest(ctx context.Context, req *messagepb.Request) (*messagepb.Response, error) {
	tracer := otel.Tracer("ocpp-go/service/message")
	ctx, span := tracer.Start(ctx, "MeterValuesRequest")
	defer span.End()

	if err := h.Service.MeterValuesRequest(ctx, req); err != nil {
		return &messagepb.Response{
			Message: "Failed to send meter values request",
		}, err
	}

	return &messagepb.Response{
		Message: "Meter values request sent",
	}, nil
}

func (h *MessageGrpcTransport) MeterValuesConfirmation(ctx context.Context, req *messagepb.Request) (*messagepb.Response, error) {
	tracer := otel.Tracer("ocpp-go/service/message")
	ctx, span := tracer.Start(ctx, "MeterValuesConfirmation")
	defer span.End()

	if err := h.Service.MeterValuesConfirmation(ctx, req); err != nil {
		return &messagepb.Response{
			Message: "Failed to send meter values confirmation",
		}, err
	}

	return &messagepb.Response{
		Message: "Meter values confirmation sent",
	}, nil
}

func (h *MessageGrpcTransport) StartTransactionRequest(ctx context.Context, req *messagepb.Request) (*messagepb.Response, error) {
	tracer := otel.Tracer("ocpp-go/service/message")
	ctx, span := tracer.Start(ctx, "StartTransactionRequest")
	defer span.End()

	if err := h.Service.StartTransactionRequest(ctx, req); err != nil {
		return &messagepb.Response{
			Message: "Failed to send start transaction request",
		}, err
	}

	return &messagepb.Response{
		Message: "Start transaction request sent",
	}, nil
}

func (h *MessageGrpcTransport) StartTransactionConfirmation(ctx context.Context, req *messagepb.Request) (*messagepb.Response, error) {
	tracer := otel.Tracer("ocpp-go/service/message")
	ctx, span := tracer.Start(ctx, "StartTransactionConfirmation")
	defer span.End()

	if err := h.Service.StartTransactionConfirmation(ctx, req); err != nil {
		return &messagepb.Response{
			Message: "Failed to send start transaction confirmation",
		}, err
	}

	return &messagepb.Response{
		Message: "Start transaction confirmation sent",
	}, nil
}

func (h *MessageGrpcTransport) StatusNotificationRequest(ctx context.Context, req *messagepb.Request) (*messagepb.Response, error) {
	tracer := otel.Tracer("ocpp-go/service/message")
	ctx, span := tracer.Start(ctx, "StatusNotificationRequest")
	defer span.End()

	if err := h.Service.StatusNotificationRequest(ctx, req); err != nil {
		return &messagepb.Response{
			Message: "Failed to send status notification request",
		}, err
	}

	return &messagepb.Response{
		Message: "Status notification request sent",
	}, nil
}

func (h *MessageGrpcTransport) StatusNotificationConfirmation(ctx context.Context, req *messagepb.Request) (*messagepb.Response, error) {
	tracer := otel.Tracer("ocpp-go/service/message")
	ctx, span := tracer.Start(ctx, "StatusNotificationConfirmation")
	defer span.End()

	if err := h.Service.StatusNotificationConfirmation(ctx, req); err != nil {
		return &messagepb.Response{
			Message: "Failed to send status notification confirmation",
		}, err
	}

	return &messagepb.Response{
		Message: "Status notification confirmation sent",
	}, nil
}

func (h *MessageGrpcTransport) StopTransactionRequest(ctx context.Context, req *messagepb.Request) (*messagepb.Response, error) {
	tracer := otel.Tracer("ocpp-go/service/message")
	ctx, span := tracer.Start(ctx, "StopTransactionRequest")
	defer span.End()

	if err := h.Service.StopTransactionRequest(ctx, req); err != nil {
		return &messagepb.Response{
			Message: "Failed to send stop transaction request",
		}, err
	}

	return &messagepb.Response{
		Message: "Stop transaction request sent",
	}, nil
}

func (h *MessageGrpcTransport) StopTransactionConfirmation(ctx context.Context, req *messagepb.Request) (*messagepb.Response, error) {
	tracer := otel.Tracer("ocpp-go/service/message")
	ctx, span := tracer.Start(ctx, "StopTransactionConfirmation")
	defer span.End()

	if err := h.Service.StopTransactionConfirmation(ctx, req); err != nil {
		return &messagepb.Response{
			Message: "Failed to send stop transaction confirmation",
		}, err
	}

	return &messagepb.Response{
		Message: "Stop transaction confirmation sent",
	}, nil
}
