package messages

import (
	"github.com/squishmeist/ocpp-go/service/ocpp/types"
)

type AuthorizeRequest struct {
	IdTag string `json:"idTag" validate:"required,max=20"`
}

type AuthorizeConfirmation struct {
	IdTagInfo *types.IdTagInfo `json:"idTagInfo" validate:"required"`
}
