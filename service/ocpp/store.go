package ocpp

import (
	"context"
	"log/slog"

	"github.com/squishmeist/ocpp-go/internal/core"
	"github.com/squishmeist/ocpp-go/service/ocpp/db/schemas"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Store struct {
	Tracer  trace.Tracer
	queries *schemas.Queries
}

func NewStore(tp trace.TracerProvider, queries *schemas.Queries) *Store {
	return &Store{
		Tracer:  tp.Tracer("store"),
		queries: queries,
	}
}

func (s *Store) GetRequestFromUuid(ctx context.Context, payload *string) (schemas.Message, core.HandlerResponse) {
	ctx, span := core.TraceDB(ctx, s.Tracer, "Store.GetRequestFromUuid")
	defer span.End()

	message, err := s.queries.GetRequestMessageByUuid(ctx, *payload)
	if err != nil {
		return schemas.Message{}, handleDBError(ctx, "to create chargepoint", err)
	}

	return message, core.HandlerResponse{
		Message: "request found",
		TraceID: trace.SpanContextFromContext(ctx).TraceID().String(),
	}
}

func (s *Store) AddRequestMessage(ctx context.Context, payload *schemas.InsertMessageParams) (schemas.Message, core.HandlerResponse) {
	ctx, span := core.TraceDB(ctx, s.Tracer, "Store.AddRequestMessage")
	defer span.End()

	response, err := s.queries.InsertMessage(ctx, *payload)
	if err != nil {
		return schemas.Message{}, handleDBError(ctx, "to add request message", err)
	}

	return response, core.HandlerResponse{
		Message: "request added",
		TraceID: trace.SpanContextFromContext(ctx).TraceID().String(),
	}
}

func (s *Store) AddConfirmationMessage(ctx context.Context, payload *schemas.InsertMessageParams) (schemas.Message, core.HandlerResponse) {
	ctx, span := core.TraceDB(ctx, s.Tracer, "Store.AddConfirmationMessage")
	defer span.End()

	// Check for existing REQUEST with the same uuid
	message, err := s.queries.GetRequestMessageByUuid(ctx, payload.Uuid)
	if err != nil {
		return schemas.Message{}, handleDBError(ctx, "no matching REQUEST for confirmation", err)
	}

	response, err := s.queries.InsertMessage(ctx, schemas.InsertMessageParams{
		Uuid:    payload.Uuid,
		Type:    "CONFIRMATION",
		Action:  message.Action,
		Payload: payload.Payload,
	})
	if err != nil {
		return schemas.Message{}, handleDBError(ctx, "to add confirmation message", err)
	}

	return response, core.HandlerResponse{
		Message: "confirmation added",
		TraceID: trace.SpanContextFromContext(ctx).TraceID().String(),
	}
}

func (s *Store) CreateChargepoint(ctx context.Context, payload *schemas.InsertChargePointParams) (schemas.Chargepoint, core.HandlerResponse) {
	ctx, span := core.TraceDB(ctx, s.Tracer, "Store.CreateChargepoint")
	defer span.End()

	chargepoint, err := s.queries.InsertChargePoint(ctx, *payload)
	if err != nil {
		return schemas.Chargepoint{}, handleDBError(ctx, "to create chargepoint", err)
	}

	return chargepoint, core.HandlerResponse{
		Message: "chargepoint created",
		TraceID: trace.SpanContextFromContext(ctx).TraceID().String(),
	}
}

func handleDBError(ctx context.Context, operation string, err error) core.HandlerResponse {
	slog.Error("failed "+operation, "error", err)
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, "failed "+operation)

	errMsg := err.Error()
	return core.HandlerResponse{
		Error:   &errMsg,
		Message: "failed " + operation,
		TraceID: trace.SpanContextFromContext(ctx).TraceID().String(),
	}
}
