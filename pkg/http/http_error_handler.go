package sharedhttp

import (
	"encoding/json"
	"net/http"

	"github.com/Macaquit0/Tropical-BFF/pkg/errors"
	"github.com/google/uuid"
)

func (r *Router) ErrorHandler(w http.ResponseWriter, rr *http.Request, err error) {
	r.log.Error(rr.Context()).Msg("error on process request: %s", err.Error())
	var genericError any

	tracerId := rr.Header.Get("x-request-id")

	if tracerId == "" {
		tracerId = uuid.NewString()
	}

	var statusCode int
	switch err := err.(type) {
	case *errors.InternalServerError:
		_, e := errors.AsInternalServerError(err)
		statusCode = http.StatusInternalServerError
		genericError = e
		break
	case *errors.ValidationError:
		_, e := errors.AsValidationError(err)
		statusCode = http.StatusBadRequest
		genericError = e
		break
	case *errors.DuplicatedEntryError:
		_, e := errors.AsDuplicatedEntryError(err)
		statusCode = http.StatusConflict
		genericError = e
		break
	case *errors.UnauthorizedError:
		_, e := errors.AsUnauthorizedError(err)
		statusCode = http.StatusUnauthorized
		genericError = e
		break
	case *errors.NotFoundError:
		_, e := errors.AsNotFoundError(err)
		statusCode = http.StatusNotFound
		genericError = e
		break
	case *errors.ForbiddenError:
		_, e := errors.AsForbiddenErrorError(err)
		statusCode = http.StatusForbidden
		genericError = e
		break
	default:
		statusCode = http.StatusInternalServerError
		genericError = errors.NewInternalServerErrorWithError("internal server error", err)
		break
	}

	body, _ := json.Marshal(genericError)

	w.WriteHeader(statusCode)
	w.Write(body)
}
