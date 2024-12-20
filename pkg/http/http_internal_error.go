package sharedhttp

type HttpInternalError struct {
	Message   string `json:"message"`
	TracerId  string `json:"tracer_id"`
	ErrorCode string `json:"error_code"`
}
