package gqlclient_test

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"

	gql "github.com/weavedev/go-gqlclient"
)

type SuiteJSONRequestBuilder struct {
	suite.Suite
}

func TestSuiteJSONRequestBuilder(t *testing.T) {
	s := SuiteJSONRequestBuilder{}
	suite.Run(t, &s)
}

func (s *SuiteJSONRequestBuilder) TestEndpoint() {
	req := gql.NewRequest("query {}")
	r, err := gql.JSONRequestBuilder("https://endpoint/query", req)
	s.NoError(err)
	s.Equal(http.MethodPost, r.Method)
	s.Equal("https", r.URL.Scheme)
	s.Equal("endpoint", r.URL.Host)
	s.Equal("/query", r.URL.Path)
}

func (s *SuiteJSONRequestBuilder) TestInvalidEndpoint() {
	req := gql.NewRequest("query {}")
	_, err := gql.JSONRequestBuilder("\r", req)
	var urlErr *url.Error
	s.ErrorAs(err, &urlErr)
}

func (s *SuiteJSONRequestBuilder) TestBody() {
	req := gql.NewRequest("query {}", gql.WithVar("key", "value"))
	r, err := gql.JSONRequestBuilder("https://endpoint/query", req)
	s.NoError(err)

	bodyBuffer := new(bytes.Buffer)
	_, err = bodyBuffer.ReadFrom(r.Body)
	s.NoError(err)
	s.Equal(`{"query":"query {}","variables":{"key":"value"}}`+"\n", bodyBuffer.String())
}

func (s *SuiteJSONRequestBuilder) TestContentType() {
	req := gql.NewRequest("query {}")
	r, err := gql.JSONRequestBuilder("https://endpoint/query", req)
	s.NoError(err)

	s.Equal("application/json; charset=utf-8", r.Header.Get("Content-Type"))
}
