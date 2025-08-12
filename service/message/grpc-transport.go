package message

import (
	"context"

	messagepb "github.com/squishmeist/ocpp-go/pkg/api/proto/message/v1"
)

type MessageServiceInterface interface {
	BootNotificationRequest(context.Context, *messagepb.Request) error
	BootNotificationConfirmation(context.Context, *messagepb.Request) error
	HeartbeatRequest(context.Context, *messagepb.Request) error
	HeartbeatConfirmation(context.Context, *messagepb.Request) error
	StatusNotificationRequest(context.Context, *messagepb.Request) error
	StatusNotificationConfirmation(context.Context, *messagepb.Request) error
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
	if err := h.Service.BootNotificationConfirmation(ctx, req); err != nil {
		return &messagepb.Response{
			Message: "Failed to send boot notification confirmation",
		}, err
	}

	return &messagepb.Response{
		Message: "Boot notification confirmation sent",
	}, nil
}

func (h *MessageGrpcTransport) StatusNotificationRequest(ctx context.Context, req *messagepb.Request) (*messagepb.Response, error) {
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
	if err := h.Service.StatusNotificationConfirmation(ctx, req); err != nil {
		return &messagepb.Response{
			Message: "Failed to send status notification confirmation",
		}, err
	}

	return &messagepb.Response{
		Message: "Status notification confirmation sent",
	}, nil
}

func (h *MessageGrpcTransport) HeartbeatRequest(ctx context.Context, req *messagepb.Request) (*messagepb.Response, error) {
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
	if err := h.Service.HeartbeatConfirmation(ctx, req); err != nil {
		return &messagepb.Response{
			Message: "Failed to send heartbeat confirmation",
		}, err
	}

	return &messagepb.Response{
		Message: "Heartbeat confirmation sent",
	}, nil
}
