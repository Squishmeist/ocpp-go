package types

type RemoteStartTransactionRequest struct {
	ConnectorId     *int             `json:"connectorId,omitempty" validate:"omitempty,gt=0"`
	IdTag           string           `json:"idTag" validate:"required,max=20"`
	ChargingProfile *ChargingProfile `json:"chargingProfile,omitempty"`
}

type RemoteStartTransactionConfirmation struct {
	Status RemoteStartStopStatus `json:"status" validate:"required,remoteStartStopStatus16"`
}
