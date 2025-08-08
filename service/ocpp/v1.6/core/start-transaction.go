package core

import "github.com/squishmeist/ocpp-go/service/ocpp/v1.6/types"

// -------------------- Start Transaction (CP -> CS) --------------------
// The Charge Point SHALL send a StartTransactionRequest to the Central System to inform about a transaction that has been started.
// If this transaction ends a reservation (see ReserveNow operation), then the StartTransaction MUST contain the reservationId.
// Upon receipt of a StartTransactionRequest, the Central System SHOULD respond with a StartTransactionConfirmation.
// This response payload MUST include a transaction id and an authorization status value.
// The Central System MUST verify validity of the identifier in the StartTransactionRequest, because the identifier might have been authorized locally by the Charge Point using outdated information.
// If Charge Point has implemented an Authorization Cache, then upon receipt of a StartTransactionConfirmation the Charge Point SHALL update the cache entry, if the idTag is not in the Local Authorization List, with the IdTagInfo value from the response as described under Authorization Cache.
// It is likely that The Central System applies sanity checks to the data contained in a StartTransactionRequest it received.
// The outcome of such sanity checks SHOULD NOT ever cause the Central System to not respond with a StartTransactionConfirmation.
// Failing to respond with a StartTransactionConfirmation will only cause the Charge Point to try the same message again as specified in Error responses to transaction-related messages.
const StartTransaction = "StartTransaction"

// This field definition of the StartTransactionRequest payload sent by the Charge Point to the Central System.
type StartTransactionRequest struct {
	ConnectorId   int             `json:"connectorId" validate:"gt=0"`
	IdTag         string          `json:"idTag" validate:"required,max=20"`
	MeterStart    int             `json:"meterStart" validate:"gte=0"`
	ReservationId *int            `json:"reservationId,omitempty" validate:"omitempty"`
	Timestamp     *types.DateTime `json:"timestamp" validate:"required"`
}

// This field definition of the StartTransactionConfirmation payload sent by the Central System to the Charge Point in response to a StartTransactionRequest.
type StartTransactionConfirmation struct {
	IdTagInfo     *types.IdTagInfo `json:"idTagInfo" validate:"required"`
	TransactionId int              `json:"transactionId"`
}
