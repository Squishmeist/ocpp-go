package ocpp

import (
	"context"

	v16 "github.com/squishmeist/ocpp-go/service/ocpp/v1.6"
	"github.com/squishmeist/ocpp-go/service/ocpp/v1.6/core"
)

type StoreAdapter interface {
	AddChargepoint(ctx context.Context, payload core.BootNotificationRequest) error
	UpdateLastHeartbeat(ctx context.Context, serialnumber string, payload core.HeartbeatConfirmation) error
}

type CacheAdapter interface {
	HasProcessed(ctx context.Context, id string) (bool, error)
	AddProcessed(ctx context.Context, id string) error
	GetRequestFromUuid(ctx context.Context, uuid string) (v16.RequestBody, error)
	AddRequest(ctx context.Context, meta v16.Meta, request v16.RequestBody) error
	RemoveRequest(ctx context.Context, meta v16.Meta, request v16.ConfirmationBody) error
}
