package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	// ErrInternal HTTP 500
	ErrInternal = &Error{
		Code:    http.StatusInternalServerError,
		Message: "Something went wrong",
	}
	// ErrUnprocessableEntity HTTP 422
	ErrUnprocessableEntity = &Error{
		Code:    http.StatusUnprocessableEntity,
		Message: "Unprocessable Entity",
	}
	// ErrBadRequest HTTP 400
	ErrBadRequest = &Error{
		Code:    http.StatusBadRequest,
		Message: "Error invalid argument",
	}
	// ErrPostNotFound HTTP 404
	ErrPostNotFound = &Error{
		Code:    http.StatusNotFound,
		Message: "Post not found",
	}
	// ErrObjectIsRequired HTTP 400
	ErrObjectIsRequired = &Error{
		Code:    http.StatusBadRequest,
		Message: "Request object should be provided",
	}
	// ErrValidPostIDIsRequired HTTP 400
	ErrValidPostIDIsRequired = &Error{
		Code:    http.StatusBadRequest,
		Message: "A valid Post id is required",
	}
	// ErrPostTimingIsRequired HTTP 400
	ErrPostTimingIsRequired = &Error{
		Code:    http.StatusBadRequest,
		Message: "Post start time and end time should be provided",
	}
	// ErrInvalidLimit HTTP 400
	ErrInvalidLimit = &Error{
		Code:    http.StatusBadRequest,
		Message: "Limit should be an integral value",
	}
	// ErrInvalidTimeFormat HTTP 400
	ErrInvalidTimeFormat = &Error{
		Code:    http.StatusBadRequest,
		Message: "Time Should be passed in RFC3339 Format: " + time.RFC3339,
	}
)

// Error main object for error
type Error struct {
	Code    int
	Message string
}

func (err *Error) Error() string {
	return err.String()
}

func (err *Error) String() string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("error: code=%s message=%s", http.StatusText(err.Code), err.Message)
}

// JSON convert Error in json
func (err *Error) JSON() []byte {
	if err == nil {
		return []byte("{}")
	}
	res, _ := json.Marshal(err)
	return res
}

// StatusCode get status code
func (err *Error) StatusCode() int {
	if err == nil {
		return http.StatusOK
	}
	return err.Code
}
