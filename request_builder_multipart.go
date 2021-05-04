package gqlclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)

// MultipartRequestBuilder creates an http.Request based on a GraphQL Request using multipart encoding.
func MultipartRequestBuilder(endpoint string, req *Request) (*http.Request, error) {
	// Encode the request as multipart request
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	if err := writer.WriteField("query", req.Query); err != nil {
		return nil, fmt.Errorf("write query field: %w", err)
	}

	// Encode and add the variables to the multipart request body.
	if len(req.Variables) > 0 {
		variablesField, err := writer.CreateFormField("variables")
		if err != nil {
			return nil, fmt.Errorf("create variables field: %w", err)
		}
		if err := json.NewEncoder(variablesField).Encode(req.Variables); err != nil {
			return nil, fmt.Errorf("encode variables: %w", err)
		}
	}

	// Close the multipart.Writer to finish the requestBody buffer.
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close writer: %w", err)
	}

	// Create a http POST request with the multipart body
	r, err := http.NewRequest(http.MethodPost, endpoint, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("create multipart request: %w", err)
	}

	// Set multipart content type
	r.Header.Set("Content-Type", writer.FormDataContentType())

	return r, nil
}
