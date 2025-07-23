package ocpp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/squishmeist/ocpp-go/internal/core/util"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func handleConfirmationBody(ctx context.Context, body ConfirmationBody, state *State) error {
	tracer := otel.Tracer("ocpp-receiver")
	ctx, span := tracer.Start(ctx, "handleConfirmationBody", trace.WithAttributes(
		attribute.String("uuid", body.Uuid),
		attribute.String("type", string(body.Type)),
	))
	defer span.End()

	errored := func(msg string, err error) error {
		return util.JustErrWithSpan2(span, msg, err)
	}

	match, err := state.FindByUuid(body.Uuid)
	if err != nil {
		return errored("RequestBody:", err)
	}
	if match.Confirmation != nil {
		return errored("confirmation already exists for uuid", fmt.Errorf("confirmation already exists for uuid: %s", body.Uuid))
	}

	span.SetAttributes(
		attribute.String("action", string(match.Request.Action)),
	)

	switch match.Request.Action {
	case Heartbeat:
		if err := heartbeatConfirmation(ctx, body.Payload); err != nil {
			return errored("Failed to handle Heartbeat confirmation", err)
		}
	case BootNotification:
		if err := bootnotificationConfirmation(ctx, body.Payload); err != nil {
			return errored("Failed to handle BootNotification confirmation", err)
		}
	default:
		slog.Error("Unknown action", "action", match.Request.Action)
		return errored("Unknown action", fmt.Errorf("unknown action: %s", match.Request.Action))
	}

	state.AddConfirmation(ConfirmationBody{
		Type:    body.Type,
		Uuid:    body.Uuid,
		Payload: body.Payload,
	})
	return nil
}

func heartbeatConfirmation(ctx context.Context, payload []byte) error {
	tracer := otel.Tracer("ocpp-receiver.heartbeatRequest")
	_, span := tracer.Start(ctx, "heartbeatRequest", trace.WithAttributes(
		attribute.String("payload", string(payload)),
	))
	defer span.End()

	obj, err := util.UnmarshalAndValidate[core.HeartbeatConfirmation](payload)
	if err != nil {
		return util.JustErrWithSpan2(span, "Failed to unmarshal HeartbeatConfirmation", err)
	}

	slog.Debug("HeartbeatConfirmation", "confirmation", obj)
	return nil
}

func bootnotificationConfirmation(ctx context.Context, payload []byte) error {
	tracer := otel.Tracer("ocpp-receiver.bootNotification")
	_, span := tracer.Start(ctx, "bootNotification", trace.WithAttributes(
		attribute.String("payload", string(payload)),
	))
	defer span.End()

	obj, err := util.UnmarshalAndValidate[core.BootNotificationConfirmation](payload)
	if err != nil {
		return util.JustErrWithSpan2(span, "Failed to unmarshal BootNotificationConfirmation", err)
	}

	slog.Debug("BootNotificationConfirmation", "confirmation", obj)
	return nil
}
