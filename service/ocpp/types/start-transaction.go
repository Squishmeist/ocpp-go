package types

type StartTransactionRequest struct {
	ConnectorId   int       `json:"connectorId" validate:"gt=0"`
	IdTag         string    `json:"idTag" validate:"required,max=20"`
	MeterStart    int       `json:"meterStart" validate:"gte=0"`
	ReservationId *int      `json:"reservationId,omitempty" validate:"omitempty"`
	Timestamp     *DateTime `json:"timestamp" validate:"required"`
}

type StartTransactionConfirmation struct {
	IdTagInfo     *IdTagInfo `json:"idTagInfo" validate:"required"`
	TransactionId int        `json:"transactionId"`
}
