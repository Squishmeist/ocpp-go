package types

type AuthorizeRequest struct {
	IdTag string `json:"idTag" validate:"required,max=20"`
}

type AuthorizeConfirmation struct {
	IdTagInfo *IdTagInfo `json:"idTagInfo" validate:"required"`
}
