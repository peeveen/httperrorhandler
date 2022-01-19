package httperrorhandler

import (
	"net/http"
)

// HTTPError describes an error that occurred during the handling of an HTTP request.
type HTTPError struct {
	// The status code that you reckon fits the error best.
	HTTPStatus int
	// The error that has occurred.
	Error error
	// Extra info you want to provide, in case "error" was produced by a third party, and you
	// don't want to fill your API code with hundreds of lines like
	// fmt.Sprintf("Failed to do something because %s", err)
	Message string
}

// HTTPErrorHandler is a function that handles an HTTP error, probably by writing the error details to the response.
type HTTPErrorHandler func(w http.ResponseWriter, r *http.Request, httpError *HTTPError)

// HandleErrors makes a call to the given hander function, and, in the event of an HTTP server error result, calls your
// supplied error handler function.
func HandleErrors(w http.ResponseWriter, r *http.Request, handlerFunc func(w http.ResponseWriter, r *http.Request) *HTTPError, errorHandler HTTPErrorHandler) {
	err := handlerFunc(w, r)
	// If we get back an error from the handlerFunc, write it to the response and set the appropriate status.
	if err != nil {
		errorHandler(w, r, err)
	}
}
