package errors

import (
	"fmt"
)

type ValidationError struct {
	GenericError
	Errors []struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	} `json:"errors"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation server error: %s", e.Message)
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		GenericError: GenericError{Message: message},
	}
}

func (e *ValidationError) AddError(field, message string) *ValidationError {
	e.Errors = append(e.Errors, struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}{Field: field, Message: message})
	return e
}

func AsValidationError(err error) (bool, ValidationError) {
	e, ok := err.(*ValidationError)
	if !ok {
		return false, ValidationError{}
	}
	return ok, *e
}
