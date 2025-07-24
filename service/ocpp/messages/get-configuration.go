package messages

type ConfigurationKey struct {
	Key      string  `json:"key" validate:"required,max=50"`
	Readonly bool    `json:"readonly"`
	Value    *string `json:"value,omitempty" validate:"omitempty,max=500"`
}

type GetConfigurationRequest struct {
	Key []string `json:"key,omitempty" validate:"omitempty,unique,dive,max=50"`
}

type GetConfigurationConfirmation struct {
	ConfigurationKey []ConfigurationKey `json:"configurationKey,omitempty" validate:"omitempty,dive"`
	UnknownKey       []string           `json:"unknownKey,omitempty" validate:"omitempty,dive,max=50"`
}
