package errors

import (
	"net/http"
)

type Error struct {
	err    error
	status int
}

func New(err error, status int) error {
	e := &Error{
		err:    err,
		status: status,
	}
	return e
}

func (e *Error) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return e.Message()
}

func (e *Error) Message() string {
	switch e.status {
	case http.StatusBadRequest:
		return "The request was invalid"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusForbidden:
		return "You do not have permission to perform this action"
	case http.StatusNotFound:
		return "The requested resource was not found"
	default:
		return "Something went wrong"
	}
}

func (e *Error) Status() int {
	if e.status >= http.StatusBadRequest {
		return e.status
	}

	return http.StatusInternalServerError
}

func (e *Error) Unwrap() error {
	return e.err
}
