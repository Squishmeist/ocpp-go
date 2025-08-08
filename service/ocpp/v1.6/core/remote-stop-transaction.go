package core

import "github.com/squishmeist/ocpp-go/service/ocpp/v1.6/types"

// -------------------- Remote Stop Transaction (CS -> CP) --------------------
// Central System can request a Charge Point to stop a transaction by sending a RemoteStopTransactionRequest to Charge Point with the identifier of the transaction.
// Charge Point SHALL reply with RemoteStopTransactionConfirmation and a status indicating whether it has accepted the request and a transaction with the given transactionId is ongoing and will be stopped.
// This remote request to stop a transaction is equal to a local action to stop a transaction.
// Therefore, the transaction SHALL be stopped, The Charge Point SHALL send a StopTransactionRequest and, if applicable, unlock the connector.
// The following two main use cases are the reason for Remote Stop Transaction:
// • Enable a CPO operator to help an EV driver that has problems stopping a transaction.
// • Enable mobile apps to control charging transactions via the Central System.
const RemoteStopTransaction = "RemoteStopTransaction"

// The field definition of the RemoteStopTransaction request payload sent by the Central System to the Charge Point.
type RemoteStopTransactionRequest struct {
	TransactionId int `json:"transactionId"`
}

// This field definition of the RemoteStopTransaction confirmation payload, sent by the Charge Point to the Central System in response to a RemoteStopTransactionRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type RemoteStopTransactionConfirmation struct {
	Status types.RemoteStartStopStatus `json:"status" validate:"required,remoteStartStopStatus16"`
}
