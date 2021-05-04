package gqlclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// JSONRequestBuilder creates an http.Request based on a GraphQL Request using a json encoding.
func JSONRequestBuilder(endpoint string, req *Request) (*http.Request, error) {
	// Encode the request as json
	var requestBody bytes.Buffer
	if err := json.NewEncoder(&requestBody).Encode(req); err != nil {
		return nil, fmt.Errorf("encode request body as json: %w", err)
	}

	// Create a http POST request with the json body
	r, err := http.NewRequest(http.MethodPost, endpoint, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("create json request: %w", err)
	}

	// Set json content type
	r.Header.Set("Content-Type", "application/json; charset=utf-8")

	return r, nil
}
