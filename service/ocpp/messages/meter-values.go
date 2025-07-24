package messages

import "github.com/squishmeist/ocpp-go/service/ocpp/types"

type MeterValuesRequest struct {
	ConnectorId   int                `json:"connectorId" validate:"gte=0"`
	TransactionId *int               `json:"transactionId,omitempty"`
	MeterValue    []types.MeterValue `json:"meterValue" validate:"required,min=1,dive"`
}

type MeterValuesConfirmation struct {
}
