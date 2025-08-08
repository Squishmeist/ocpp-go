package ocpp

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	iCore "github.com/squishmeist/ocpp-go/internal/core"
	"github.com/squishmeist/ocpp-go/service/ocpp/db/schemas"
	"github.com/squishmeist/ocpp-go/service/ocpp/v1.6/core"
	"github.com/squishmeist/ocpp-go/service/ocpp/v1.6/types"
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

func (s *DbStore) AddChargepoint(ctx context.Context, payload core.BootNotificationRequest) error {
	ctx, span := iCore.TraceDB(ctx, s.Tracer, "Store.AddChargepoint")
	defer span.End()

	_, err := s.queries.InsertChargepoint(ctx, schemas.InsertChargepointParams{
		SerialNumber:      payload.ChargeBoxSerialNumber,
		Model:             payload.ChargePointModel,
		Vendor:            payload.ChargePointVendor,
		FirmwareVersion:   payload.FirmwareVersion,
		Iicid:             sql.NullString{String: payload.Iccid, Valid: payload.Iccid != ""},
		Imsi:              sql.NullString{String: payload.Imsi, Valid: payload.Imsi != ""},
		MeterSerialNumber: sql.NullString{String: payload.MeterSerialNumber, Valid: payload.MeterSerialNumber != ""},
		MeterType:         sql.NullString{String: payload.MeterType, Valid: payload.MeterType != ""},
		LastBoot:          types.Now().Time,
	})
	if err != nil {
		return handleDBError(ctx, "to add chargepoint", err)
	}

	return nil
}

func (s *DbStore) UpdateLastHeartbeat(ctx context.Context, serialnumber string, payload core.HeartbeatConfirmation) error {
	ctx, span := iCore.TraceDB(ctx, s.Tracer, "Store.UpdateLastHeartbeat")
	defer span.End()

	result, err := s.queries.UpdateChargepointLastHeartbeat(ctx, schemas.UpdateChargepointLastHeartbeatParams{
		SerialNumber:  serialnumber,
		LastHeartbeat: sql.NullTime{Time: payload.CurrentTime.Time, Valid: true},
	})
	if err != nil {
		return handleDBError(ctx, "to update last heartbeat", err)
	}
	if result == "" {
		return fmt.Errorf("serial number %s not found", serialnumber)
	}

	return nil
}

func handleDBError(ctx context.Context, operation string, err error) error {
	slog.Error("failed "+operation, "error", err)
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, "failed "+operation)

	return fmt.Errorf("failed %s: %w", operation, err)
}
