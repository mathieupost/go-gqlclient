// Package gqlclient provides a client for interacting with a GraphQL endpoint.
//
// Create a client (safe to share across requests)
//  client := gqlclient.NewClient(
//      "https://localhost/graphql",
//      // Optionally supply options:
//      // Set default headers.
//      gqlclient.WithDefaultHeader("Authorization", "Bearer " + token),
//      // Use a custom http.Client.
//      gqlclient.WithHTTPClient(customClient),
//      // Use another request builder (default: gqlclient.JSONRequestBuilder).
//      gqlclient.WithRequestBuilder(gqlclient.MultipartRequestBuilder),
//  )
//
// Make a request
//  req := gqlclient.NewRequest(`
//      query ($key: String!) {
//          item(id: $key) {
//              field1
//              field2
//              field3
//          }
//      }`,
//      // Optionally supply options:
//      // Set any variables.
//      gqlclient.WithVar("key", "value"),
//      // Set header fields.
//      gqlclient.WithHeader("Cache-Control", "no-cache"),
//      // Pass a Context for the request (default: context.Background()).
//      gqlclient.WithContext(ctx),
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
//
// Inspect the returned GraphQL errors
//  var gqlerrs gqlclient.ErrorList
//  if errors.As(err, &gqlerrs) {
//      // Check path
//      println(gqlerrs[0].Path)
//  }
package gqlclient
