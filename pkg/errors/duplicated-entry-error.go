package errors

import (
	"fmt"
)

type DuplicatedEntryError struct {
	GenericError
}

func (e *DuplicatedEntryError) Error() string {
	return fmt.Sprintf("duplicated entry error: %s", e.Message)
}

func NewDuplicatedEntryError(message string) *DuplicatedEntryError {
	return &DuplicatedEntryError{
		GenericError{
			Message: message,
		},
	}
}

func AsDuplicatedEntryError(err error) (bool, DuplicatedEntryError) {
	e, ok := err.(*DuplicatedEntryError)
	if !ok {
		return false, DuplicatedEntryError{}
	}
	return ok, *e
}

func IsDuplicatedEntryError(err error) bool {
	_, ok := err.(*DuplicatedEntryError)
	if !ok {
		return false
	}
	return ok
}
