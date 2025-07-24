package messages

import "github.com/squishmeist/ocpp-go/service/ocpp/types"

type RemoteStartTransactionRequest struct {
	ConnectorId     *int                   `json:"connectorId,omitempty" validate:"omitempty,gt=0"`
	IdTag           string                 `json:"idTag" validate:"required,max=20"`
	ChargingProfile *types.ChargingProfile `json:"chargingProfile,omitempty"`
}

type RemoteStartTransactionConfirmation struct {
	Status types.RemoteStartStopStatus `json:"status" validate:"required,remoteStartStopStatus16"`
}
