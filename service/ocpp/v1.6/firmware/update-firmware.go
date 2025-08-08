package firmware

import "github.com/squishmeist/ocpp-go/service/ocpp/v1.6/types"

// -------------------- Get Diagnostics (CS -> CP) --------------------
// Central System can notify a Charge Point that it needs to update its firmware.
// The Central System SHALL send an UpdateFirmwareRequest to instruct the Charge Point to install new firmware.
// The payload SHALL contain a date and time after which the Charge Point is allowed to retrieve the new firmware and the location from which the firmware can be downloaded.
// The Charge Point SHALL respond with a UpdateFirmwareConfirmation. The Charge Point SHOULD start retrieving the firmware as soon as possible after retrieve-date.
// During downloading and installation of the firmware, the Charge Point MUST send FirmwareStatusNotificationRequest payloads to keep the Central System updated with the status of the update process.
// The Charge Point SHALL, if the new firmware image is "valid", install the new firmware as soon as it is able to.
// If it is not possible to continue charging during installation of firmware, it is RECOMMENDED to wait until Charging Session has ended (Charge Point idle) before commencing installation.
// It is RECOMMENDED to set connectors that are not in use to UNAVAILABLE while the Charge Point waits for the Session to end.
const UpdateFirmware = "UpdateFirmware"

// The field definition of the UpdateFirmware request payload sent by the Central System to the Charge Point.
type UpdateFirmwareRequest struct {
	Location      string          `json:"location" validate:"required,uri"`
	Retries       *int            `json:"retries,omitempty" validate:"omitempty,gte=0"`
	RetrieveDate  *types.DateTime `json:"retrieveDate" validate:"required"`
	RetryInterval *int            `json:"retryInterval,omitempty" validate:"omitempty,gte=0"`
}

// This field definition of the UpdateFirmware confirmation payload, sent by the Charge Point to the Central System in response to a UpdateFirmwareRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type UpdateFirmwareConfirmation struct {
}
