package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/backend/bff-cognito/pkg/errors"
	"github.com/backend/bff-cognito/pkg/logger"
	"github.com/backend/bff-cognito/pkg/trace"
)

type Rest struct {
	log        *logger.Logger
	HttpMethod HttpMethod
	Url        string
}

func New(log *logger.Logger, httpMethod HttpMethod, url string) Rest {
	return Rest{
		log: log, HttpMethod: httpMethod, Url: url,
	}
}

type HttpMethod string

const (
	GET    = HttpMethod("GET")
	POST   = HttpMethod("POST")
	PUT    = HttpMethod("PUT")
	PATCH  = HttpMethod("PATCH")
	DELETE = HttpMethod("DELETE")
)

func (r *Rest) Do(ctx context.Context, headers map[string]string, body any, res any, params ...string) error {
	for _, v := range params {
		r.Url = strings.Replace(r.Url, "{}", v, 1)
	}

	j, err := json.Marshal(body)
	if err != nil {
		return errors.NewInternalServerError("error on parse request body to json")
	}

	b := strings.NewReader(string(j))

	req, err := http.NewRequest(string(r.HttpMethod), r.Url, b)
	if err != nil {
		return err
	}
	traceId, ok := ctx.Value(trace.TraceId).(string)
	if ok {
		req.Header.Set("trace-id", traceId)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, h := range headers {
		req.Header.Set(k, h)
	}

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return errors.NewInternalServerErrorWithError(fmt.Sprintf("error do http request on %s %s", r.HttpMethod, r.Url), err)
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			r.log.Error(ctx).Msg("error on close http do request body")
		}
	}()
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.NewInternalServerErrorWithError(fmt.Sprintf("error on read response body of do request on %s %s", r.HttpMethod, r.Url), err)
	}
	switch {
	case response.StatusCode >= 200 && response.StatusCode < 300:
		if res != nil {
			err = json.Unmarshal(resBody, res)
			if err != nil {
				return errors.NewInternalServerErrorWithError(fmt.Sprintf("error on parse response body of do request on %s %s", r.HttpMethod, r.Url), err)
			}
		}

		return nil

	case response.StatusCode >= 400 && response.StatusCode < 500:
		if response.StatusCode == 401 || response.StatusCode == 403 {
			return errors.NewUnauthorizedError()
		}
		if response.StatusCode == 404 {
			return errors.NewNotFoundError(fmt.Sprintf("error not found on process %s http request to: %s", r.HttpMethod, r.Url))
		}
		if response.StatusCode == 409 {
			return errors.NewDuplicatedEntryError(fmt.Sprintf("error duplicaetd entry on process %s http request to: %s", r.HttpMethod, r.Url))
		}
		validationError := errors.ValidationError{}
		err = json.Unmarshal(resBody, &validationError)
		if err != nil {
			return errors.NewInternalServerErrorWithError(fmt.Sprintf("error on process %s http request to: %s", r.HttpMethod, r.Url), err)
		}
		return &validationError
	case response.StatusCode >= 500 && response.StatusCode < 600:
		return errors.NewInternalServerErrorWithError(fmt.Sprintf("error on process %s http request to: %s", r.HttpMethod, r.Url), fmt.Errorf("%s", string(resBody)))
	default:
		return errors.NewInternalServerErrorWithError(fmt.Sprintf("error on process %s http request to: %s", r.HttpMethod, r.Url), fmt.Errorf("%s", string(resBody)))
	}
}
