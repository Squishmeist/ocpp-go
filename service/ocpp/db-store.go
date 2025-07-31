package ocpp

import (
	"context"
	"log/slog"

	"github.com/squishmeist/ocpp-go/internal/core"
	"github.com/squishmeist/ocpp-go/internal/core/utils"
	"github.com/squishmeist/ocpp-go/service/ocpp/db/schemas"
	"github.com/squishmeist/ocpp-go/service/ocpp/types"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type DbStore struct {
	Tracer  trace.Tracer
	queries *schemas.Queries
}

func NewDbStore(tp trace.TracerProvider, queries *schemas.Queries) *DbStore {
	return &DbStore{
		Tracer:  tp.Tracer("store"),
		queries: queries,
	}
}

func (s *DbStore) GetRequestFromUuid(ctx context.Context, uuid string) (types.RequestBody, core.HandlerResponse) {
	ctx, span := core.TraceDB(ctx, s.Tracer, "Store.GetRequestFromUuid")
	defer span.End()

	response, err := s.queries.GetRequestMessageByUuid(ctx, uuid)
	if err != nil {
		return types.RequestBody{}, handleDBError(ctx, "to create chargepoint", err)
	}
	action := types.ActionKind(response.Action)
	if !action.IsValid() {
		return types.RequestBody{}, core.HandlerResponse{
			Error:   utils.StringPtr("invalid action kind"),
			Message: "invalid action kind",
			TraceID: trace.SpanContextFromContext(ctx).TraceID().String(),
		}
	}

	return types.RequestBody{
			Uuid:    response.Uuid,
			Action:  action,
			Payload: []byte(response.Payload),
		}, core.HandlerResponse{
			Message: "request found",
			TraceID: trace.SpanContextFromContext(ctx).TraceID().String(),
		}
}

func (s *DbStore) AddRequestMessage(ctx context.Context, request types.RequestBody) (types.MessageBody, core.HandlerResponse) {
	ctx, span := core.TraceDB(ctx, s.Tracer, "Store.AddRequestMessage")
	defer span.End()

	response, err := s.queries.InsertMessage(ctx, schemas.InsertMessageParams{
		Uuid:    request.Uuid,
		Type:    "REQUEST",
		Action:  string(request.Action),
		Payload: string(request.Payload),
	})
	if err != nil {
		return types.MessageBody{}, handleDBError(ctx, "to add request message", err)
	}

	action := types.ActionKind(response.Action)
	if !action.IsValid() {
		return types.MessageBody{}, core.HandlerResponse{
			Error:   utils.StringPtr("invalid action kind"),
			Message: "invalid action kind",
			TraceID: trace.SpanContextFromContext(ctx).TraceID().String(),
		}
	}

	return types.MessageBody{
			Kind:    types.Confirmation,
			Uuid:    response.Uuid,
			Action:  action,
			Payload: []byte(response.Payload),
		}, core.HandlerResponse{
			Message: "request added",
			TraceID: trace.SpanContextFromContext(ctx).TraceID().String(),
		}
}

func (s *DbStore) AddConfirmationMessage(ctx context.Context, payload types.ConfirmationBody) (types.MessageBody, core.HandlerResponse) {
	ctx, span := core.TraceDB(ctx, s.Tracer, "Store.AddConfirmationMessage")
	defer span.End()

	// Check for existing REQUEST with the same uuid
	requestMsg, err := s.queries.GetRequestMessageByUuid(ctx, payload.Uuid)
	if err != nil {
		return types.MessageBody{}, handleDBError(ctx, "no matching REQUEST for confirmation", err)
	}

	response, err := s.queries.InsertMessage(ctx, schemas.InsertMessageParams{
		Uuid:    payload.Uuid,
		Type:    "CONFIRMATION",
		Action:  requestMsg.Action,
		Payload: string(payload.Payload),
	})
	if err != nil {
		return types.MessageBody{}, handleDBError(ctx, "to add confirmation message", err)
	}
	action := types.ActionKind(response.Action)
	if !action.IsValid() {
		return types.MessageBody{}, core.HandlerResponse{
			Error:   utils.StringPtr("invalid action kind"),
			Message: "invalid action kind",
			TraceID: trace.SpanContextFromContext(ctx).TraceID().String(),
		}
	}

	return types.MessageBody{
			Kind:    types.Confirmation,
			Uuid:    response.Uuid,
			Action:  action,
			Payload: []byte(response.Payload),
		}, core.HandlerResponse{
			Message: "confirmation added",
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
