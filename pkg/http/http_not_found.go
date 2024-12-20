package sharedhttp

type HttpNotFound struct {
	Message   string `json:"message"`
	Resource  string `json:"resource"`
	TracerID  string `json:"tracer_id"`
	ErrorCode string `json:"error_code"`
}

func NewHttpNotFound(message, resource, tracerID, errorCode string) *HttpNotFound {
	return &HttpNotFound{
		Message:   message,
		ErrorCode: errorCode,
		Resource:  resource,
		TracerID:  tracerID,
	}
}
