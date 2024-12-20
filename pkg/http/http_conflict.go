package sharedhttp

type HttpConflict struct {
	Message   string            `json:"message"`
	Resource  string            `json:"resource"`
	TracerID  string            `json:"tracer_id"`
	ErrorCode string            `json:"error_code"`
	Metadata  map[string]string `json:"metadata"`
}
