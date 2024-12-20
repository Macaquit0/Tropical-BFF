package errors

import "fmt"

type NotFoundError struct {
	GenericError
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found error: %s", e.Message)
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		GenericError{
			Message: message,
		},
	}
}

func AsNotFoundError(err error) (bool, NotFoundError) {
	e, ok := err.(*NotFoundError)
	if !ok {
		return false, NotFoundError{}
	}
	return ok, *e
}

func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	if !ok {
		return false
	}
	return ok
}
