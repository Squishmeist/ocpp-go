package ocpp

import (
	"context"
	"log/slog"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/squishmeist/ocpp-go/internal/core/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type RequestHandlerProps struct {
	ctx    context.Context
	body   RequestBody
	tracer trace.Tracer
}

type RequestHandler interface {
	Handle(RequestHandlerProps) error
}

type HeartbeatRequestHandler struct{}

func (h HeartbeatRequestHandler) Handle(props RequestHandlerProps) error {
	ctx, body, tracer := props.ctx, props.body, props.tracer

	_, span := tracer.Start(ctx, string(Heartbeat), trace.WithAttributes(
		attribute.String("uuid", body.Uuid),
		attribute.String("type", string(body.Type)),
		attribute.String("action", string(body.Action)),
		attribute.String("payload", string(body.Payload)),
	))
	defer span.End()

	obj, err := util.UnmarshalAndValidate[core.HeartbeatRequest](body.Payload)
	if err != nil {
		return util.JustErrWithSpan(span, "Failed to unmarshal HeartbeatRequest", err)
	}

	slog.Debug("HeartbeatRequest", "request", obj)
	// add the request to the state
	state.AddRequest(body)
	return nil
}

type BootNotificationRequestHandler struct{}

func (h BootNotificationRequestHandler) Handle(props RequestHandlerProps) error {
	ctx, body, tracer := props.ctx, props.body, props.tracer

	_, span := tracer.Start(ctx, string(BootNotification), trace.WithAttributes(
		attribute.String("uuid", body.Uuid),
		attribute.String("type", string(body.Type)),
		attribute.String("action", string(body.Action)),
		attribute.String("payload", string(body.Payload)),
	))
	defer span.End()

	obj, err := util.UnmarshalAndValidate[core.BootNotificationRequest](body.Payload)
	if err != nil {
		return util.JustErrWithSpan(span, "Failed to unmarshal BootNotificationRequest", err)
	}

	slog.Debug("BootNotificationRequest", "request", obj)
	return nil

}

var requestHandlers = map[ActionType]RequestHandler{
	Heartbeat:        HeartbeatRequestHandler{},
	BootNotification: BootNotificationRequestHandler{},
}

func getRequestHandler(action ActionType) (RequestHandler, bool) {
	handler, ok := requestHandlers[action]
	return handler, ok
}
