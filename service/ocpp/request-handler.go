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

func handleRequestBody(ctx context.Context, body RequestBody, state *State) error {
	tracer := otel.Tracer("ocpp-receiver")
	ctx, span := tracer.Start(ctx, "handleRequestBody")
	span.SetAttributes(
		attribute.String("uuid", body.Uuid),
		attribute.String("type", string(body.Type)),
		attribute.String("action", string(body.Action)),
	)

	switch body.Action {
	case Heartbeat:
		if err := heartbeatRequest(ctx, body.Payload); err != nil {
			message := "Failed to handle Heartbeat request"
			span.RecordError(err)
			span.SetStatus(codes.Error, message)
			span.End()
			return fmt.Errorf("%s: %w", message, err)
		}
	case BootNotification:
		if err := bootnotificationRequest(ctx, body.Payload); err != nil {
			message := "Failed to handle BootNotification request"
			span.RecordError(err)
			span.SetStatus(codes.Error, message)
			span.End()
			return fmt.Errorf("%s: %w", message, err)
		}
	default:
		slog.Error("Unknown action", "action", body.Action)
		span.SetStatus(codes.Error, "Unknown action")
		span.End()
		return fmt.Errorf("unknown action: %s", body.Action)
	}

	state.AddRequest(RequestBody{Type: body.Type, Uuid: body.Uuid, Action: body.Action})
	span.End()
	return nil
}

func heartbeatRequest(ctx context.Context, payload []byte) error {
	tracer := otel.Tracer("ocpp-receiver.heartbeatRequest")
	_, span := tracer.Start(ctx, "heartbeatRequest")
	span.SetAttributes(
		attribute.String("payload", string(payload)),
	)

	obj, err := unmarshalAndValidate[core.HeartbeatRequest](payload)
	if err != nil {
		message := "Failed to unmarshal HeartbeatRequest"
		span.RecordError(err)
		span.SetStatus(codes.Error, message)
		span.End()
		return fmt.Errorf("%s: %w", message, err)
	}

	slog.Debug("HeartbeatRequest", "request", obj)
	span.End()
	return nil
}

func bootnotificationRequest(ctx context.Context, payload []byte) error {
	tracer := otel.Tracer("ocpp-receiver.bootnotificationRequest")
	_, span := tracer.Start(ctx, "bootnotificationRequest")
	span.SetAttributes(
		attribute.String("payload", string(payload)),
	)

	obj, err := unmarshalAndValidate[core.BootNotificationRequest](payload)
	if err != nil {
		message := "Failed to unmarshal BootNotificationRequest"
		span.RecordError(err)
		span.SetStatus(codes.Error, message)
		span.End()
		return fmt.Errorf("%s: %w", message, err)
	}
	slog.Debug("BootNotificationRequest", "request", obj)

	span.End()
	return nil
}
