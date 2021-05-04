// Package gqlclient provides a client for interacting with a GraphQL endpoint.
//
//  import (
//      gql "github.com/weavedev/go-gqlclient"
//      "github.com/weavedev/go-gqlclient/builders"
//  )
//
// Create a client (safe to share across requests)
//  client := gql.NewClient(
//      "https://localhost/graphql",
//      // Optionally supply options:
//      // Set default headers.
//      gql.WithDefaultHeader("Authorization", "Bearer " + token),
//      // Use a custom http.Client.
//      gql.WithHTTPClient(customClient),
//      // Use another request builder (default: builders.JSON).
//      gql.WithRequestBuilder(builders.Multipart),
//  )
//
// Make a request
//  req := gql.NewRequest(`
//      query ($key: String!) {
//        item(id: $key) {
//          field1
//          field2
//          field3
//        }
//      }`,
//      // Optionally supply options:
//      // Set any variables.
//      gql.WithVar("key", "value"),
//      // Set header fields.
//      gql.WithHeader("Cache-Control", "no-cache"),
//      // Pass a Context for the request (default: context.Background()).
//      gql.WithContext(ctx),
//  )
//
// Do the request and capture the response.
//  var resp struct {
//      Item struct {
//          Field1 string
//          Field2 string
//          Field3 string
//      }
//  }
//  err := client.Do(req, &resp)
package gqlclient
