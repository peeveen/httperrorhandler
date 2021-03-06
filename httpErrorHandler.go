package httperrorhandler

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// Error describes an error that occurred during the handling of an HTTP request.
type Error struct {
	// The error that has occurred (if this HTTP error was generated by an internal error).
	cause error
	// The status code that you reckon fits the error best.
	Status int `json:"status,omitempty"`
	// Human-readable message that describes the error.
	Detail string `json:"detail,omitempty"`
	// Optional URI identifier the error instance.
	Instance string `json:"instance,omitempty"`
	// URI identifying the type of the error.
	Type string `json:"type"`
}

// Cause returns the error that was the cause of this error.
func (e *Error) Cause() error {
	return e.cause
}

// Error returns the message that describes this error.
func (e *Error) Error() string {
	return e.Detail
}

func chooseNonEmptyString(str1 string, str2 string) string {
	if str1 == "" {
		return str2
	}
	return str1
}

// Wrap returns an HTTP error object that wraps the given error.
func Wrap(e error, httpError *Error) *Error {
	return &Error{cause: e, Status: httpError.Status, Detail: chooseNonEmptyString(httpError.Detail, e.Error()), Instance: httpError.Instance, Type: httpError.Type}
}

// Handler is a function that handles an HTTP error, probably by writing the error details to the response.
type Handler func(w http.ResponseWriter, r *http.Request, httpError *Error)

// Handle makes a call to the given hander function, and, in the event of an HTTP server error result, calls your
// supplied error handler function.
func Handle(w http.ResponseWriter, r *http.Request, handlerFunc func(w http.ResponseWriter, r *http.Request) *Error, errorHandlerFunc Handler) {
	err := handlerFunc(w, r)
	// If we get back an error from the handlerFunc, write it to the response and set the appropriate status.
	if err != nil && errorHandlerFunc != nil {
		errorHandlerFunc(w, r, err)
	}
}

// DefaultErrorHandler is a default handler for HTTP errors that you can use if you choose.
// It will write the HTTP status and error JSON to the response. If the object cannot be
// marshalled to JSON, it will write a plain text response of the error Detail field.
func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, e *Error) {
	logrus.Error(e.Detail)
	w.WriteHeader(e.Status)
	jsonError, err := json.Marshal(e)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.Write([]byte(e.Detail))
	} else {
		w.Header().Set("Content-Type", "application/problem+json; charset=UTF-8")
		w.Write(jsonError)
	}
}
