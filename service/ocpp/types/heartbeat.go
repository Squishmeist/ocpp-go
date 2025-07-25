package types

import (
	"github.com/go-playground/validator/v10"
)

type HeartbeatRequest struct {
}

type HeartbeatConfirmation struct {
	CurrentTime *DateTime `json:"currentTime" validate:"required"`
}

func validateHeartbeatConfirmation(sl validator.StructLevel) {
	confirmation := sl.Current().Interface().(HeartbeatConfirmation)
	if DateTimeIsNull(confirmation.CurrentTime) {
		sl.ReportError(confirmation.CurrentTime, "CurrentTime", "currentTime", "required", "")
	}
}

func init() {
	Validate.RegisterStructValidation(validateHeartbeatConfirmation, HeartbeatConfirmation{})
}
