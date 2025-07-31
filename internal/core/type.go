package core

type HandlerResponse struct {
	Error   *string `json:"error"`
	Message string  `json:"message"`
	TraceID string  `json:"trace_id"`
}
