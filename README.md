# gqlclient [![GoDoc](https://godoc.org/github.com/weavedev/go-gqlclient?status.png)](http://godoc.org/github.com/weavedev/go-gqlclient)

GraphQL client for Go.

* Simple, familiar API
* Respects `context.Context` timeouts and cancellation
* Build and execute a GraphQL request using json or multipart
* Use strong Go types for response data
* Use variables, custom headers and a custom http client
* Advanced error handling

## Installation

Make sure you have a working Go environment. To install gqlclient, simply run:

```
$ go get github.com/weavedev/go-gqlclient
```

## Usage

```go
import (
	gql "github.com/weavedev/go-gqlclient"
)


// Create a client (safe to share across requests)
client := gql.NewClient(
	"https://localhost/graphql",
	// Optionally supply options:
	// Set default headers.
	gql.WithDefaultHeader("Authorization", "Bearer " + token),
	// Use a custom http.Client.
	gql.WithHTTPClient(customClient),
	// Use another request builder (default: gql.JSONRequestBuilder).
	gql.WithRequestBuilder(gql.MultipartRequestBuilder),
)

// Make a request
req := gql.NewRequest(`
	query ($key: String!) {
		item(id: $key) {
			field1
			field2
			field3
		}
	}`,
	// Optionally supply options:
	// Set any variables.
	gql.WithVar("key", "value"),
	// Set header fields.
	gql.WithHeader("Cache-Control", "no-cache"),
	// Pass a Context for the request (default: context.Background()).
	gql.WithContext(ctx),
)

// Do the request and capture the response.
var resp struct {
	Item struct {
		Field1 string
		Field2 string
		Field3 string
	}
}
err := client.Do(req, &resp)

// Inspect the returned GraphQL errors
var gqlerrs gql.ErrorList
if errors.As(err, &gqlerrs) {
	// Check path
	println(gqlerrs[0].Path)
}
```

## Thanks

Inspired by https://github.com/machinebox/graphql
