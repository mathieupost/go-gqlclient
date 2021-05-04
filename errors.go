package gqlclient

import (
	"errors"
	"fmt"
	"net/http"
)

// HTTPError represents an error that occurred in the http transport layer and not in the GraphQL layer.
type HTTPError struct {
	StatusCode int
}

func (e *HTTPError) Error() string {
	statusText := http.StatusText(e.StatusCode)
	if statusText == "" {
		return fmt.Sprintf("http error: %d", e.StatusCode)
	} else {
		return fmt.Sprintf("http error: %d %s", e.StatusCode, statusText)
	}
}

// NewHTTPError	creates a new HTTPError with the given http status code.
func NewHTTPError(statusCode int) *HTTPError {
	return &HTTPError{StatusCode: statusCode}
}

// ErrBadResponse is used when the response body cannot be parsed.
var ErrBadResponse = errors.New("response was not GraphQL compliant")
