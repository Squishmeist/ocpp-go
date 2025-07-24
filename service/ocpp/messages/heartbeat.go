package messages

import (
	"github.com/go-playground/validator/v10"
	"github.com/squishmeist/ocpp-go/service/ocpp/types"
)

type HeartbeatRequest struct {
}

type HeartbeatConfirmation struct {
	CurrentTime *types.DateTime `json:"currentTime" validate:"required"`
}

func validateHeartbeatConfirmation(sl validator.StructLevel) {
	confirmation := sl.Current().Interface().(HeartbeatConfirmation)
	if types.DateTimeIsNull(confirmation.CurrentTime) {
		sl.ReportError(confirmation.CurrentTime, "CurrentTime", "currentTime", "required", "")
	}
}

func init() {
	types.Validate.RegisterStructValidation(validateHeartbeatConfirmation, HeartbeatConfirmation{})
}
