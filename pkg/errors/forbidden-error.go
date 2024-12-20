package errors

import "fmt"

type ForbiddenError struct {
	GenericError
}

func (e *ForbiddenError) Error() string {
	return fmt.Sprintf("forbidden error: %s", e.Message)
}

func NewForbiddenErrorError() *ForbiddenError {
	return &ForbiddenError{
		GenericError{
			Message: "unauthorized",
		},
	}
}

func AsForbiddenErrorError(err error) (bool, ForbiddenError) {
	e, ok := err.(*ForbiddenError)
	if !ok {
		return false, ForbiddenError{}
	}
	return ok, *e
}

func IsForbiddenErrorError(err error) bool {
	_, ok := err.(*ForbiddenError)
	if !ok {
		return false
	}
	return ok
}
