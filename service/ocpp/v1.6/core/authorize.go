package core

import "github.com/squishmeist/ocpp-go/service/ocpp/v1.6/types"

// -------------------- Authorize (CP -> CS) --------------------
// Before the owner of an electric vehicle can start or stop charging, the Charge Point has to authorize the operation.
// Upon receipt of an AuthorizeRequest, the Central System SHALL respond with an AuthorizeConfirmation.
// This response payload SHALL indicate whether or not the idTag is accepted by the Central System.
// If the Central System accepts the idTag then the response payload MAY include a parentIdTag and MUST include an authorization status value indicating acceptance or a reason for rejection.
// A Charge Point MAY authorize identifier locally without involving the Central System, as described in Local Authorization List.
// The Charge Point SHALL only supply energy after authorization.
const Authorize = "Authorize"

// The field definition of the Authorize request payload sent by the Charge Point to the Central System.
type AuthorizeRequest struct {
	IdTag string `json:"idTag" validate:"required,max=20"`
}

// This field definition of the Authorize confirmation payload, sent by the Charge Point to the Central System in response to an AuthorizeRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type AuthorizeConfirmation struct {
	IdTagInfo *types.IdTagInfo `json:"idTagInfo" validate:"required"`
}
