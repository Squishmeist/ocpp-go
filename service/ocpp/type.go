package ocpp

type OcppMessageType string

const (
	HeartbeatConfirmation        OcppMessageType = "HeartbeatConfirmation"
	HeartbeatRequest             OcppMessageType = "HeartbeatRequest"
	BootNotificationRequest      OcppMessageType = "BootNotificationRequest"
	BootNotificationConfirmation OcppMessageType = "BootNotificationConfirmation"
	Unknown                      OcppMessageType = "Unknown"
)

type OcppMessage struct {
	Type OcppMessageType
	Data any
}

type Body struct {
	ChargePointID string `json:"chargePointId"`
	Payload       string `json:"payload"`
}