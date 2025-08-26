package client

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHttpClient(t *testing.T) {
	// Test creating a new HTTP client
	httpClient := &http.Client{}
	bearerToken := "test-bearer-token"

	client := NewHttpClient(httpClient, bearerToken)

	assert.NotNil(t, client)

	// Verify it's the correct type
	clientImpl, ok := client.(*HttpClientImpl)
	assert.True(t, ok)
	assert.Equal(t, httpClient, clientImpl.client)
	assert.Equal(t, bearerToken, clientImpl.bearerToken)
}

func TestNewHttpClient_WithEmptyToken(t *testing.T) {
	// Test creating a client with empty token
	httpClient := &http.Client{}
	bearerToken := ""

	client := NewHttpClient(httpClient, bearerToken)

	assert.NotNil(t, client)

	clientImpl, ok := client.(*HttpClientImpl)
	assert.True(t, ok)
	assert.Equal(t, httpClient, clientImpl.client)
	assert.Empty(t, clientImpl.bearerToken)
}

func TestNewHttpClient_WithNilHttpClient(t *testing.T) {
	// Test creating a client with nil http.Client
	bearerToken := "test-token"

	client := NewHttpClient(nil, bearerToken)

	assert.NotNil(t, client)

	clientImpl, ok := client.(*HttpClientImpl)
	assert.True(t, ok)
	assert.Nil(t, clientImpl.client)
	assert.Equal(t, bearerToken, clientImpl.bearerToken)
}

func TestHttpClientImpl_Structure(t *testing.T) {
	// Test the structure of HttpClientImpl
	httpClient := &http.Client{}
	bearerToken := "test-structure-token"

	client := NewHttpClient(httpClient, bearerToken)
	clientImpl := client.(*HttpClientImpl)

	// Verify fields are accessible and correct
	assert.Equal(t, httpClient, clientImpl.client)
	assert.Equal(t, bearerToken, clientImpl.bearerToken)
	assert.NotEmpty(t, clientImpl.bearerToken)
}

func TestHttpClientImpl_Post(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			body, err := ioutil.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, `{"key":"value"}`, string(body))

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		}))
		defer server.Close()

		client := NewHttpClient(server.Client(), "test-token")
		payload := map[string]string{"key": "value"}

		resp, err := client.Post(context.Background(), payload, server.URL)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"status":"ok"}`, string(respBody))
	})

	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewHttpClient(server.Client(), "test-token")
		payload := map[string]string{"key": "value"}

		resp, err := client.Post(context.Background(), payload, server.URL)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("invalid url", func(t *testing.T) {
		client := NewHttpClient(&http.Client{}, "test-token")
		payload := map[string]string{"key": "value"}

		_, err := client.Post(context.Background(), payload, "invalid-url")

		assert.Error(t, err)
	})

	t.Run("payload marshal error", func(t *testing.T) {
		client := NewHttpClient(&http.Client{}, "test-token")
		payload := make(chan int) // Invalid payload for JSON marshaling

		_, err := client.Post(context.Background(), payload, "http://localhost")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to marshal payload")
	})
}
