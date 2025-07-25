package types

type RemoteStopTransactionRequest struct {
	TransactionId int `json:"transactionId"`
}

type RemoteStopTransactionConfirmation struct {
	Status RemoteStartStopStatus `json:"status" validate:"required,remoteStartStopStatus16"`
}
