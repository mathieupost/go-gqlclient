package gqlclient

import "context"

// Request is a GraphQL request.
type Request struct {
	Context   context.Context        `json:"-"`
	Headers   map[string]string      `json:"-"`
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// NewRequest makes a new Request with the specified string.
func NewRequest(query string, opts ...RequestOption) *Request {
	req := &Request{
		Context:   context.Background(),
		Headers:   make(map[string]string),
		Query:     query,
		Variables: make(map[string]interface{}),
	}
	for _, optionFunc := range opts {
		optionFunc(req)
	}
	return req
}

// RequestOption are functions that are passed into NewRequest to modify the Request.
type RequestOption func(*Request)

// WithContext sets the Context which is used when executing the Request.
//  NewRequest(query, WithContext(ctx))
func WithContext(ctx context.Context) RequestOption {
	return func(r *Request) {
		r.Context = ctx
	}
}

// WithHeader sets an entry in the header of a Request to the specified value.
//  NewRequest(query, WithHeader(key, value))
func WithHeader(key string, value string) RequestOption {
	return func(r *Request) {
		r.Headers[key] = value
	}
}

// WithVar defines the value of a variable in the query of a Request.
//  NewRequest(query, WithVar(name, value))
func WithVar(name string, value interface{}) RequestOption {
	return func(r *Request) {
		r.Variables[name] = value
	}
}
