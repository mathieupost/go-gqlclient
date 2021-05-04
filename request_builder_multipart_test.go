package gqlclient_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"

	gql "github.com/weavedev/go-gqlclient"
)

type SuiteMultipart struct {
	suite.Suite
}

func TestSuiteMultipart(t *testing.T) {
	s := SuiteMultipart{}
	suite.Run(t, &s)
}

func (s *SuiteMultipart) TestEndpoint() {
	req := gql.NewRequest("query {}")
	r, err := gql.Multipart("https://endpoint/query", req)
	s.NoError(err)
	s.Equal(http.MethodPost, r.Method)
	s.Equal("https", r.URL.Scheme)
	s.Equal("endpoint", r.URL.Host)
	s.Equal("/query", r.URL.Path)
}

func (s *SuiteMultipart) TestInvalidEndpoint() {
	req := gql.NewRequest("query {}")
	_, err := gql.Multipart("\r", req)
	var urlErr *url.Error
	s.ErrorAs(err, &urlErr)
}

func (s *SuiteMultipart) TestBody() {
	req := gql.NewRequest("query {}", gql.WithVar("key", "value"))
	r, err := gql.Multipart("https://endpoint/query", req)
	s.NoError(err)

	s.Equal("query {}", r.PostFormValue("query"))
	s.Equal(`{"key":"value"}`+"\n", r.PostFormValue("variables"))
}

func (s *SuiteMultipart) TestContentType() {
	req := gql.NewRequest("query {}")
	r, err := gql.Multipart("https://endpoint/query", req)
	s.NoError(err)

	s.Contains(r.Header.Get("Content-Type"), "multipart/form-data;")
}
