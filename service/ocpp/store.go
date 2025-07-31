package ocpp

import (
	"context"

	"github.com/squishmeist/ocpp-go/internal/core"
	"github.com/squishmeist/ocpp-go/service/ocpp/types"
)

type StoreAdapter interface {
	GetRequestFromUuid(ctx context.Context, uuid string) (types.RequestBody, core.HandlerResponse)
	AddRequestMessage(ctx context.Context, request types.RequestBody) (types.MessageBody, core.HandlerResponse)
	AddConfirmationMessage(ctx context.Context, confirmation types.ConfirmationBody) (types.MessageBody, core.HandlerResponse)
}
