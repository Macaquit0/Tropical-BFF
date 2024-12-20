package errors

import "errors"

type GenericError struct {
	Message string `json:"message"`
}

func Is(err error, targe error) bool {
	return errors.Is(err, targe)
}

func As(err error, targe any) bool {
	return errors.As(err, targe)
}
