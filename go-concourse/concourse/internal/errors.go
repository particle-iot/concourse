package internal

import (
	"errors"
	"fmt"

	"github.com/google/jsonapi"
)

type UnexpectedResponseError struct {
	error
	StatusCode int
	Status     string
	Body       string
}

func (e UnexpectedResponseError) Error() string {
	return fmt.Sprintf("Unexpected Response\nStatus: %s\nBody: %s", e.Status, e.Body)
}

type ForbiddenError struct {
	error
	Reason string
}

func (e ForbiddenError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("forbidden: %s", e.Reason)
	}
	return "forbidden"
}

type ResourceNotFoundError jsonapi.ErrorsPayload

func (e ResourceNotFoundError) Error() string {
	if len(e.Errors) == 0 {
		return "resource not found"
	} else {
		var response string

		for i, error := range e.Errors {
			if i > 0 {
				response = response + " "
			}
			response = response + error.Detail
		}

		return response
	}
}

var ErrUnauthorized = errors.New("not authorized")
var ErrForbidden = ForbiddenError{}
