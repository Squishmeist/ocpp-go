package firmware

import (
	"github.com/go-playground/validator/v10"
	"github.com/squishmeist/ocpp-go/service/ocpp/v1.6/types"
)

// -------------------- Firmware Status Notification (CP -> CS) --------------------
// The Charge Point sends a notification to inform the Central System about the progress of the downloading and installation of a firmware update.
// The Charge Point SHALL only send the status Idle after receipt of a TriggerMessage for a Firmware Status Notification, when it is not busy downloading/installing firmware.
// The FirmwareStatusNotification requests SHALL be sent to keep the Central System updated with the status of the update process.
const FirmwareStatusNotification = "FirmwareStatusNotification"

// Status reported in FirmwareStatusNotificationRequest.
type FirmwareStatus string

const (
	FirmwareStatusDownloaded         FirmwareStatus = "Downloaded"
	FirmwareStatusDownloadFailed     FirmwareStatus = "DownloadFailed"
	FirmwareStatusDownloading        FirmwareStatus = "Downloading"
	FirmwareStatusIdle               FirmwareStatus = "Idle"
	FirmwareStatusInstallationFailed FirmwareStatus = "InstallationFailed"
	FirmwareStatusInstalling         FirmwareStatus = "Installing"
	FirmwareStatusInstalled          FirmwareStatus = "Installed"
)

func isValidFirmwareStatus(fl validator.FieldLevel) bool {
	status := FirmwareStatus(fl.Field().String())
	switch status {
	case FirmwareStatusDownloaded, FirmwareStatusDownloadFailed, FirmwareStatusDownloading, FirmwareStatusIdle, FirmwareStatusInstallationFailed, FirmwareStatusInstalling, FirmwareStatusInstalled:
		return true
	default:
		return false
	}
}

// The field definition of the FirmwareStatusNotification request payload sent by the Charge Point to the Central System.
type FirmwareStatusNotificationRequest struct {
	Status FirmwareStatus `json:"status" validate:"required,firmwareStatus16"`
}

// This field definition of the FirmwareStatusNotification confirmation payload, sent by the Central System to the Charge Point in response to a FirmwareStatusNotificationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type FirmwareStatusNotificationConfirmation struct {
}

func init() {
	_ = types.Validate.RegisterValidation("firmwareStatus16", isValidFirmwareStatus)
}
