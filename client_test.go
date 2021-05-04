package gqlclient_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	gql "github.com/weavedev/go-gqlclient"
	"github.com/weavedev/go-gqlclient/mocks"
)

type SuiteClient struct {
	suite.Suite
}

func TestSuiteClient(t *testing.T) {
	s := SuiteClient{}
	suite.Run(t, &s)
}

func (s *SuiteClient) TestDefault() {
	c := gql.NewClient("test")

	s.Equal("test", c.Endpoint)
	s.Equal(map[string]string{"Accept": "application/json; charset=utf-8"}, c.DefaultHeaders)
	s.Equal(http.DefaultClient, c.HTTPClient)
}

func (s *SuiteClient) TestClientOption() {
	clientOption := new(mocks.ClientOption)
	clientOption.
		On("Execute", mock.AnythingOfType("*gqlclient.Client")).
		Return()

	c := gql.NewClient("test", clientOption.Execute)
	clientOption.AssertExpectations(s.T())
	clientOption.AssertCalled(s.T(), "Execute", c)
}

func (s *SuiteClient) TestDo() {
	var resp struct {
		Value string
	}

	httpClient := new(mocks.HTTPClient)
	httpClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(&http.Response{
			Body:       ioutil.NopCloser(strings.NewReader(`{"data": {"value": "some data"}}`)),
			StatusCode: http.StatusOK,
		}, nil)

	c := gql.NewClient("test", gql.WithHTTPClient(httpClient))
	err := c.Do(gql.NewRequest(""), &resp)
	httpClient.AssertExpectations(s.T())
	s.NoError(err)
	s.Equal("some data", resp.Value)
}

func (s *SuiteClient) TestSkipDecoding() {
	var resp *struct {
		Users []struct {
			Name string
		}
	}

	httpClient := new(mocks.HTTPClient)
	httpClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(&http.Response{
			Body:       ioutil.NopCloser(strings.NewReader(`{"data": {"value": "some data"}}`)),
			StatusCode: http.StatusOK,
		}, nil)

	c := gql.NewClient("test", gql.WithHTTPClient(httpClient))
	// Here resp is nil, so decoding should be skipped.
	err := c.Do(gql.NewRequest(""), resp)
	httpClient.AssertExpectations(s.T())
	s.NoError(err)
	s.Nil(resp)
}

func (s *SuiteClient) TestRequestBuilderError() {
	berr := errors.New("fail")

	requestBuilder := new(mocks.RequestBuilder)
	requestBuilder.
		On("Execute", "endpoint", mock.AnythingOfType("*gqlclient.Request")).
		Return(nil, berr)
	httpClient := gql.NewClient("endpoint", gql.WithRequestBuilder(requestBuilder.Execute))

	err := httpClient.Do(gql.NewRequest(""), nil)
	requestBuilder.AssertExpectations(s.T())
	s.ErrorIs(err, berr)
}

func (s *SuiteClient) TestWithContext() {
	ctx := context.WithValue(context.Background(), "key", "value")

	// Stub RequestBuilder, we just want to test if the correct context is set after building the http request.
	requestBuilder := func(_ string, _ *gql.Request) (*http.Request, error) {
		return http.NewRequest(http.MethodPost, "localhost", strings.NewReader(""))
	}

	httpClient := new(mocks.HTTPClient)
	httpClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Run(func(args mock.Arguments) {
			rarg := args.Get(0)
			s.Require().IsType((*http.Request)(nil), rarg)
			r := rarg.(*http.Request)
			s.Equal(ctx, r.Context())
			s.Equal("value", r.Context().Value("key"))
		}).Return(nil, errors.New("stop"))

	c := gql.NewClient("test",
		gql.WithRequestBuilder(requestBuilder),
		gql.WithHTTPClient(httpClient))
	_ = c.Do(gql.NewRequest("", gql.WithContext(ctx)), nil)
	httpClient.AssertExpectations(s.T())
}

func (s *SuiteClient) TestDefaultHeaders() {
	httpClient := new(mocks.HTTPClient)
	httpClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Run(func(args mock.Arguments) {
			rarg := args.Get(0)
			s.Require().IsType((*http.Request)(nil), rarg)
			r := rarg.(*http.Request)
			s.Equal("test-value", r.Header.Get("test-header"))
		}).
		Return(&http.Response{
			Body:       ioutil.NopCloser(strings.NewReader("{}")),
			StatusCode: http.StatusOK,
		}, nil)

	c := gql.NewClient("test", gql.WithHTTPClient(httpClient),
		gql.WithDefaultHeader("test-header", "test-value"))
	err := c.Do(gql.NewRequest(""), nil)
	httpClient.AssertExpectations(s.T())
	s.NoError(err)
}

func (s *SuiteClient) TestRequestHeaders() {
	httpClient := new(mocks.HTTPClient)
	httpClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Run(func(args mock.Arguments) {
			rarg := args.Get(0)
			s.Require().IsType((*http.Request)(nil), rarg)
			r := rarg.(*http.Request)
			s.Equal("test-value", r.Header.Get("test-header"))
		}).
		Return(&http.Response{
			Body:       ioutil.NopCloser(strings.NewReader("{}")),
			StatusCode: http.StatusOK,
		}, nil)

	c := gql.NewClient("test", gql.WithHTTPClient(httpClient))
	err := c.Do(gql.NewRequest("",
		gql.WithHeader("test-header", "test-value")), nil)
	httpClient.AssertExpectations(s.T())
	s.NoError(err)
}

func (s *SuiteClient) TestHTTPClientDoError() {
	derr := errors.New("error")
	httpClient := new(mocks.HTTPClient)
	httpClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(nil, derr)

	c := gql.NewClient("test", gql.WithHTTPClient(httpClient))
	err := c.Do(gql.NewRequest(""), nil)

	httpClient.AssertExpectations(s.T())
	s.ErrorIs(err, derr)
}

func (s *SuiteClient) TestCloseBodyError() {
	closeErr := errors.New("error on close")

	// Return error on close
	mockBody := mocks.NewBody("{}")
	mockBody.On("Close").Return(closeErr)

	mockClient := new(mocks.HTTPClient)
	mockClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(&http.Response{
			Body:       mockBody,
			StatusCode: http.StatusOK,
		}, nil)

	c := gql.NewClient("test", gql.WithHTTPClient(mockClient))
	err := c.Do(gql.NewRequest(""), nil)

	mockClient.AssertExpectations(s.T())
	s.ErrorIs(err, closeErr)
}

func (s *SuiteClient) TestHTTPError() {
	httpClient := new(mocks.HTTPClient)
	httpClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(&http.Response{
			// Return a non-200 status code
			StatusCode: http.StatusInternalServerError,
			// in combination with an empty body, to simulate an internal error.
			Body: ioutil.NopCloser(strings.NewReader("")),
		}, nil)

	c := gql.NewClient("test", gql.WithHTTPClient(httpClient))
	err := c.Do(gql.NewRequest(""), nil)

	httpClient.AssertExpectations(s.T())
	var herr *gql.HTTPError
	s.ErrorAs(err, &herr)
	s.Equal(500, herr.StatusCode)
}

func (s *SuiteClient) TestBadResponseError() {
	httpClient := new(mocks.HTTPClient)
	httpClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(&http.Response{
			// An empty string is not valid JSON and will throw an error on decoding.
			Body:       ioutil.NopCloser(strings.NewReader("")),
			StatusCode: http.StatusOK,
		}, nil)

	c := gql.NewClient("test", gql.WithHTTPClient(httpClient))
	err := c.Do(gql.NewRequest(""), nil)

	httpClient.AssertExpectations(s.T())
	s.ErrorIs(err, gql.ErrBadResponse)
}

func (s *SuiteClient) TestGQLError() {
	var resp struct {
		Value string
	}

	httpClient := new(mocks.HTTPClient)
	httpClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(&http.Response{
			// An empty string is not valid JSON and will throw an error on decoding.
			Body: ioutil.NopCloser(strings.NewReader(`{
				"data": {"value": "some data"},
				"errors": [{
					"message": "invalid query",
					"path": ["someList",1,"someField"],
					"locations": [{"line":1,"column":2}]
				}]
			}`)),
			StatusCode: http.StatusOK,
		}, nil)

	c := gql.NewClient("test", gql.WithHTTPClient(httpClient))
	err := c.Do(gql.NewRequest(""), &resp)

	httpClient.AssertExpectations(s.T())
	var herr gql.ErrorList
	s.Assert().ErrorAs(err, &herr)
	s.Assert().Len(herr, 1)
	s.Equal("invalid query", herr[0].Message)
}
