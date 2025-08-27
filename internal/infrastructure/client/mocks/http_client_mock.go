package mocks

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/stretchr/testify/mock"
)

// MockHTTPClient is a mock implementation of HTTPClient for testing
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Post(ctx context.Context, headers map[string]string, payload interface{}, url string) (*http.Response, error) {
	args := m.Called(ctx, payload, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

// CreateMockResponse Helper function to create mock HTTP responses
func CreateMockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
}
