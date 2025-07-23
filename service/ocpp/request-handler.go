package ocpp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/squishmeist/ocpp-go/internal/core/util"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func handleRequestBody(ctx context.Context, body RequestBody, state *State) error {
	tracer := otel.Tracer("ocpp-receiver")
	ctx, span := tracer.Start(ctx, "handleRequestBody", trace.WithAttributes(
		attribute.String("uuid", body.Uuid),
		attribute.String("type", string(body.Type)),
		attribute.String("action", string(body.Action)),
	))
	defer span.End()

	errored := func(msg string, err error) error {
		return util.JustErrWithSpan2(span, msg, err)
	}

	switch body.Action {
	case Heartbeat:
		if err := heartbeatRequest(ctx, body.Payload); err != nil {
			return errored("Failed to handle Heartbeat request", err)
		}
	case BootNotification:
		if err := bootnotificationRequest(ctx, body.Payload); err != nil {
			return errored("Failed to handle BootNotification request", err)
		}
	default:
		slog.Error("Unknown action", "action", body.Action)
		return errored("Unknown action", fmt.Errorf("unknown action: %s", body.Action))
	}

	state.AddRequest(RequestBody{Type: body.Type, Uuid: body.Uuid, Action: body.Action})
	return nil
}

func heartbeatRequest(ctx context.Context, payload []byte) error {
	tracer := otel.Tracer("ocpp-receiver.heartbeatRequest")
	_, span := tracer.Start(ctx, "heartbeatRequest", trace.WithAttributes(
		attribute.String("payload", string(payload)),
	))
	defer span.End()

	obj, err := util.UnmarshalAndValidate[core.HeartbeatRequest](payload)
	if err != nil {
		return util.JustErrWithSpan2(span, "Failed to unmarshal HeartbeatRequest", err)
	}

	slog.Debug("HeartbeatRequest", "request", obj)
	return nil
}

func bootnotificationRequest(ctx context.Context, payload []byte) error {
	tracer := otel.Tracer("ocpp-receiver.bootnotificationRequest")
	_, span := tracer.Start(ctx, "bootnotificationRequest", trace.WithAttributes(
		attribute.String("payload", string(payload)),
	))
	defer span.End()

	obj, err := util.UnmarshalAndValidate[core.BootNotificationRequest](payload)
	if err != nil {
		message := "Failed to unmarshal BootNotificationRequest"
		span.RecordError(err)
		span.SetStatus(codes.Error, message)
		span.End()
		return fmt.Errorf("%s: %w", message, err)
	}

	slog.Debug("BootNotificationRequest", "request", obj)
	return nil
}
