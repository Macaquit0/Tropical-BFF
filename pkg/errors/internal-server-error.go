package errors

import (
	"fmt"
)

type InternalServerError struct {
	GenericError
}

func (e *InternalServerError) Error() string {
	return fmt.Sprintf("internal server error: %s", e.Message)
}

func NewInternalServerError(message string) *InternalServerError {
	return &InternalServerError{
		GenericError{
			Message: message,
		},
	}
}

func NewInternalServerErrorWithError(message string, err error) *InternalServerError {
	return &InternalServerError{
		GenericError{
			Message: fmt.Sprintf("%s. error: %s", message, err),
		},
	}
}

func AsInternalServerError(err error) (bool, InternalServerError) {
	e, ok := err.(*InternalServerError)
	if !ok {
		return false, InternalServerError{}
	}
	return ok, *e
}

func IsInternalServerError(err error) bool {
	_, ok := err.(*InternalServerError)
	if !ok {
		return false
	}
	return ok
}
