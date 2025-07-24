package messages

import "github.com/squishmeist/ocpp-go/service/ocpp/types"

type RemoteStopTransactionRequest struct {
	TransactionId int `json:"transactionId"`
}

type RemoteStopTransactionConfirmation struct {
	Status types.RemoteStartStopStatus `json:"status" validate:"required,remoteStartStopStatus16"`
}
