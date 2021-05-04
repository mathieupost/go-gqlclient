package gqlclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type RequestBuilder func(endpoint string, req *Request) (*http.Request, error)

// Client is a client for interacting with a GraphQL API.
type Client struct {
	endpoint       string
	httpClient     HTTPClient
	defaultHeaders map[string]string
	requestBuilder RequestBuilder
}

// NewClient makes a new Client capable of making GraphQL requests.
func NewClient(endpoint string, opts ...ClientOption) *Client {
	client := &Client{
		endpoint:       endpoint,
		httpClient:     http.DefaultClient,
		defaultHeaders: make(map[string]string),
		requestBuilder: JSONRequestBuilder,
	}

	// Set default Accept header
	client.defaultHeaders["Accept"] = "application/json; charset=utf-8"

	// Parse options
	for _, optionFunc := range opts {
		optionFunc(client)
	}

	return client
}

// Do executes the Request and decodes the response from the data field into the given response
// object. Pass in a nil response object to skip response parsing. If the request fails or the
// server returns an error, the first error will be returned.
func (c *Client) Do(req *Request, resp interface{}) (err error) {
	httpReq, err := c.requestBuilder(c.endpoint, req)
	if err != nil {
		return fmt.Errorf("request builder: %w", err)
	}
	httpReq = httpReq.WithContext(req.ctx)

	// Set default headers.
	for key, value := range c.defaultHeaders {
		httpReq.Header.Set(key, value)
	}

	// Set request headers.
	for key, value := range req.headers {
		httpReq.Header.Set(key, value)
	}

	// Do the request.
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func() {
		cerr := httpResp.Body.Close()
		if cerr != nil && err == nil {
			err = fmt.Errorf("close body: %w", cerr)
		}
	}()

	// Decode the response body.
	var gqlResp responseWithErrors
	if resp == nil {
		// Skip data decoding if there is nothing to decode into. Only decode errors if they exist.
		gqlResp = &errorsResponse{}
	} else {
		gqlResp = &response{Data: resp}
	}
	if err := json.NewDecoder(httpResp.Body).Decode(gqlResp); err != nil {
		// GraphQL endpoints should always return a 200, as per GraphQL spec. So, if there was was a
		// problem decoding the response, something outside of the GraphQL layer went wrong.
		if httpResp.StatusCode != http.StatusOK {
			return NewHTTPError(httpResp.StatusCode)
		}
		return ErrBadResponse
	}

	// Return the GraphQL errors, if any.
	if len(gqlResp.getErrors()) > 0 {
		return gqlResp.getErrors()
	}
	return nil
}

// ClientOption are functions that are passed into NewClient to modify the behaviour of the Client.
type ClientOption func(*Client)

// WithHTTPClient specifies the underlying http.Client to use when making requests.
//  NewClient(endpoint, WithHTTPClient(specificHTTPClient))
func WithHTTPClient(httpclient HTTPClient) ClientOption {
	return func(client *Client) {
		client.httpClient = httpclient
	}
}

// WithDefaultHeader sets a default value for a header entry of every Request sent with this client.
//  NewClient(endpoint, WithDefaultHeader(key, value))
func WithDefaultHeader(key string, value string) ClientOption {
	return func(client *Client) {
		client.defaultHeaders[key] = value
	}
}

// WithRequestBuilder sets a function that executes the Request sent with this client.
//  NewClient(endpoint, WithDefaultHeader(key, value))
func WithRequestBuilder(builder RequestBuilder) ClientOption {
	return func(client *Client) {
		client.requestBuilder = builder
	}
}
