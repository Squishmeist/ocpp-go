package message

import (
	"context"

	ocpppb "github.com/squishmeist/ocpp-go/pkg/api/proto/ocpp/v1"
)

type MessageServiceInterface interface {
	HeartbeatRequest(context.Context, *ocpppb.Request) error
	HeartbeatConfirmation(context.Context, *ocpppb.Request) error
	BootNotificationRequest(context.Context, *ocpppb.Request) error
	BootNotificationConfirmation(context.Context, *ocpppb.Request) error
}

type MessageGrpcTransport struct {
	ocpppb.UnimplementedOCPPMessageServer
	Service MessageServiceInterface
}

func NewMessageGrpcTransport(service MessageServiceInterface) *MessageGrpcTransport {
	return &MessageGrpcTransport{
		Service: service,
	}
}

func (h *MessageGrpcTransport) HeartbeatRequest(ctx context.Context, req *ocpppb.Request) (*ocpppb.Response, error) {
	if err := h.Service.HeartbeatRequest(ctx, req); err != nil {
		return &ocpppb.Response{
			Message: "Failed to send heartbeat request",
		}, err
	}

	return &ocpppb.Response{
		Message: "Heartbeat request sent",
	}, nil
}

func (h *MessageGrpcTransport) HeartbeatConfirmation(ctx context.Context, req *ocpppb.Request) (*ocpppb.Response, error) {
	if err := h.Service.HeartbeatConfirmation(ctx, req); err != nil {
		return &ocpppb.Response{
			Message: "Failed to send heartbeat confirmation",
		}, err
	}

	return &ocpppb.Response{
		Message: "Heartbeat confirmation sent",
	}, nil
}

func (h *MessageGrpcTransport) BootNotificationRequest(ctx context.Context, req *ocpppb.Request) (*ocpppb.Response, error) {
	if err := h.Service.BootNotificationRequest(ctx, req); err != nil {
		return &ocpppb.Response{
			Message: "Failed to send boot notification request",
		}, err
	}

	return &ocpppb.Response{
		Message: "Boot notification request sent",
	}, nil
}

func (h *MessageGrpcTransport) BootNotificationConfirmation(ctx context.Context, req *ocpppb.Request) (*ocpppb.Response, error) {
	if err := h.Service.BootNotificationConfirmation(ctx, req); err != nil {
		return &ocpppb.Response{
			Message: "Failed to send boot notification confirmation",
		}, err
	}

	return &ocpppb.Response{
		Message: "Boot notification confirmation sent",
	}, nil
}
