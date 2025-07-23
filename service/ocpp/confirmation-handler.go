package ocpp

import (
	"context"
	"log/slog"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/squishmeist/ocpp-go/internal/core/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ConfirmationHandlerProps struct {
	ctx    context.Context
	body   ConfirmationBody
	tracer trace.Tracer
}

type ConfirmationHandler interface {
	Handle(ConfirmationHandlerProps) error
}

type HeartbeatConfirmationHandler struct{}

func (h HeartbeatConfirmationHandler) Handle(props ConfirmationHandlerProps) error {
	ctx, body, tracer := props.ctx, props.body, props.tracer
	_, span := tracer.Start(ctx, string(Heartbeat), trace.WithAttributes(
		attribute.String("uuid", string(body.Uuid)),
		attribute.String("type", string(body.Type)),
		attribute.String("payload", string(body.Payload)),
	))
	defer span.End()

	obj, err := util.UnmarshalAndValidate[core.HeartbeatConfirmation](body.Payload)
	if err != nil {
		return util.JustErrWithSpan(span, "Failed to unmarshal HeartbeatConfirmation", err)
	}

	slog.Debug("HeartbeatConfirmation", "confirmation", obj)
	state.AddConfirmation(ConfirmationBody{
		Type:    body.Type,
		Uuid:    body.Uuid,
		Payload: body.Payload,
	})
	return nil
}

type BootNotificationConfirmationHandler struct{}

func (h BootNotificationConfirmationHandler) Handle(props ConfirmationHandlerProps) error {
	ctx, body, tracer := props.ctx, props.body, props.tracer

	_, span := tracer.Start(ctx, string(BootNotification), trace.WithAttributes(
		attribute.String("uuid", string(body.Uuid)),
		attribute.String("type", string(body.Type)),
		attribute.String("payload", string(body.Payload)),
	))
	defer span.End()

	obj, err := util.UnmarshalAndValidate[core.BootNotificationConfirmation](body.Payload)
	if err != nil {
		return util.JustErrWithSpan(span, "Failed to unmarshal BootNotificationConfirmation", err)
	}

	slog.Debug("BootNotificationConfirmation", "confirmation", obj)
	state.AddConfirmation(ConfirmationBody{
		Type:    body.Type,
		Uuid:    body.Uuid,
		Payload: body.Payload,
	})
	return nil
}

var confirmationHandlers = map[ActionType]ConfirmationHandler{
	Heartbeat:        HeartbeatConfirmationHandler{},
	BootNotification: BootNotificationConfirmationHandler{},
}

func getConfirmationHandler(uuid string) (ConfirmationHandler, bool) {
	match, err := state.FindByUuid(uuid)
	if err != nil {
		slog.Error("Failed to find request by UUID", "uuid", uuid, "error", err)
		return nil, false
	}
	if match.Confirmation != nil {
		slog.Error("Confirmation already exists for UUID", "uuid", uuid)
		return nil, false
	}
	handler, ok := confirmationHandlers[match.Request.Action]
	return handler, ok
}
