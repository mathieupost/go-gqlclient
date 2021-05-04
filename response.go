package gqlclient

import (
	"bytes"
	"strconv"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type responseWithErrors interface {
	GetErrors() ErrorList
}

// Response contains the default data and errors entries of a GraphQL response.
type Response struct {
	Data interface{} `json:"data,omitempty"`
	ErrorsResponse
}

// ErrorsResponse contains only the errors entry of a GraphQL response.
type ErrorsResponse struct {
	Errors ErrorList `json:"errors,omitempty"`
}

// GetErrors returns the errors that were returned in a GraphQL response.
func (r ErrorsResponse) GetErrors() ErrorList {
	return r.Errors
}

// ErrorList is an error type to embed the errors list from a GraphQL response.
type ErrorList []*Error

// Error returns the first error from a GraphQL response.
func (m ErrorList) Error() string {
	if len(m) == 0 {
		return ""
	}
	return m[0].Error()
}

// Error contains all the data that a GraphQL error can contain.
type Error struct {
	Message    string                 `json:"message"`
	Path       ast.Path               `json:"path,omitempty"`
	Locations  []gqlerror.Location    `json:"locations,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// Error formats the error using locations, path and message.
func (e Error) Error() string {
	var res bytes.Buffer
	res.WriteString("graphql:")

	if len(e.Locations) > 0 {
		res.WriteString(strconv.Itoa(e.Locations[0].Line))
		res.WriteByte(':')
		res.WriteString(strconv.Itoa(e.Locations[0].Column))
		res.WriteByte(':')
	}

	if e.Path != nil {
		res.WriteByte(' ')
		res.WriteString(e.Path.String())
		res.WriteByte(':')
	}

	res.WriteByte(' ')
	res.WriteString(e.Message)

	return res.String()
}
