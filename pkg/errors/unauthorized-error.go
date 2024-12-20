package errors

import "fmt"

type UnauthorizedError struct {
	GenericError
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized error: %s", e.Message)
}

func NewUnauthorizedError() *UnauthorizedError {
	return &UnauthorizedError{
		GenericError{
			Message: "unauthorized",
		},
	}
}

func NewUnauthorizedErrorMsg(msg string) *UnauthorizedError {
	return &UnauthorizedError{
		GenericError{
			Message: msg,
		},
	}
}

func AsUnauthorizedError(err error) (bool, UnauthorizedError) {
	e, ok := err.(*UnauthorizedError)
	if !ok {
		return false, UnauthorizedError{}
	}
	return ok, *e
}

func IsUnauthorizedError(err error) bool {
	_, ok := err.(*UnauthorizedError)
	if !ok {
		return false
	}
	return ok
}
