package ocpp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func handleConfirmationBody(ctx context.Context, body ConfirmationBody, state *State) error {
	tracer := otel.Tracer("ocpp-receiver")
	ctx, span := tracer.Start(ctx, "handleConfirmationBody")
	span.SetAttributes(
		attribute.String("uuid", body.Uuid),
		attribute.String("type", string(body.Type)),
	)

	match, err := state.FindByUuid(body.Uuid)
	if err != nil {
		message := "RequestBody not found"
		span.RecordError(err)
		span.SetStatus(codes.Error, message)
		span.End()
		return fmt.Errorf("%s: %w", message, err)
	}

	if match.Confirmation != nil {
		message := "confirmation already exists for uuid"
		span.RecordError(err)
		span.SetStatus(codes.Error, message)
		span.End()
		return fmt.Errorf("%s: %w", message, err)
	}

	span.SetAttributes(
		attribute.String("action", string(match.Request.Action)),
	)

	switch match.Request.Action {
	case Heartbeat:
		if err := heartbeatConfirmation(ctx, body.Payload); err != nil {
			message := "Failed to handle Heartbeat confirmation"
			span.RecordError(err)
			span.SetStatus(codes.Error, message)
			span.End()
			return fmt.Errorf("%s: %w", message, err)
		}
	case BootNotification:
		if err := bootnotificationConfirmation(ctx, body.Payload); err != nil {
			message := "Failed to handle BootNotification confirmation"
			span.RecordError(err)
			span.SetStatus(codes.Error, message)
			span.End()
			return fmt.Errorf("%s: %w", message, err)
		}
	default:
		slog.Error("Unknown action", "action", match.Request.Action)
		span.SetStatus(codes.Error, "Unknown action")
		span.End()
		return fmt.Errorf("unknown action: %s", match.Request.Action)
	}

	state.AddConfirmation(ConfirmationBody{
		Type:    body.Type,
		Uuid:    body.Uuid,
		Payload: body.Payload,
	})

	span.End()
	return nil
}

func heartbeatConfirmation(ctx context.Context, payload []byte) error {
	tracer := otel.Tracer("ocpp-receiver.heartbeatRequest")
	_, span := tracer.Start(ctx, "heartbeatRequest")
	span.SetAttributes(
		attribute.String("payload", string(payload)),
	)

	obj, err := unmarshalAndValidate[core.HeartbeatConfirmation](payload)
	if err != nil {
		message := "Failed to unmarshal HeartbeatConfirmation"
		span.RecordError(err)
		span.SetStatus(codes.Error, message)
		span.End()
		return fmt.Errorf("%s: %w", message, err)
	}

	slog.Debug("HeartbeatConfirmation", "confirmation", obj)
	span.End()
	return nil
}

func bootnotificationConfirmation(ctx context.Context, payload []byte) error {
	tracer := otel.Tracer("ocpp-receiver.bootNotification")
	_, span := tracer.Start(ctx, "bootNotification")
	span.SetAttributes(
		attribute.String("payload", string(payload)),
	)

	obj, err := unmarshalAndValidate[core.BootNotificationConfirmation](payload)
	if err != nil {
		message := "Failed to unmarshal BootNotificationConfirmation"
		span.RecordError(err)
		span.SetStatus(codes.Error, message)
		span.End()
		return fmt.Errorf("%s: %w", message, err)
	}

	slog.Debug("BootNotificationConfirmation", "confirmation", obj)
	span.End()
	return nil
}
