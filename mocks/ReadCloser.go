package mocks

import (
	"io"
	"strings"

	mock "github.com/stretchr/testify/mock"
)

// Body is a mock type for the ReadCloser type
type Body struct {
	mock.Mock
	Reader io.Reader
}

func NewBody(content string) *Body {
	body := new(Body)
	body.Reader = strings.NewReader(content)
	return body
}

// Close provides a mock function with given fields:
func (_m *Body) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Read wraps the inner io.Reader
func (_m *Body) Read(p []byte) (int, error) {
	return _m.Reader.Read(p)
}
