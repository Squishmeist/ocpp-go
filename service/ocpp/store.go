package ocpp

import (
	"context"

	"github.com/squishmeist/ocpp-go/service/ocpp/types"
)

type StoreAdapter interface {
	AddChargepoint(ctx context.Context, payload types.BootNotificationRequest) error
	UpdateLastHeartbeat(ctx context.Context, serialnumber string, payload types.HeartbeatConfirmation) error
}

type CacheAdapter interface {
	HasProcessed(ctx context.Context, id string) (bool, error)
	GetRequestFromUuid(ctx context.Context, uuid string) (types.RequestBody, error)
	AddRequest(ctx context.Context, meta types.Meta, request types.RequestBody) error
	RemoveRequest(ctx context.Context, meta types.Meta, request types.ConfirmationBody) error
}
