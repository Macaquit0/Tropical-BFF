package sharedhttp

type HttpBadRequest struct {
	Message   string `json:"message"`
	TracerID  string `json:"tracer_id"`
	ErrorCode string `json:"error_code"`
}
