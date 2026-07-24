package weather

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientGetCurrent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/New York", r.URL.Path)
		assert.Equal(t, "j1", r.URL.Query().Get("format"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"current_condition":[{"temp_C":"25.5","weatherDesc":[{"value":"Sunny"}]}]}`)
	}))
	defer server.Close()

	client, err := NewClient(Config{BaseURL: server.URL, Timeout: time.Second})
	require.NoError(t, err)

	current, err := client.GetCurrent(context.Background(), " New York ")
	require.NoError(t, err)
	assert.Equal(t, "Sunny", current.Description)
	assert.Equal(t, 25.5, current.TemperatureC)
}

func TestClientGetCurrentRejectsUpstreamError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	client, err := NewClient(Config{BaseURL: server.URL, Timeout: time.Second})
	require.NoError(t, err)

	_, err = client.GetCurrent(context.Background(), "Beijing")
	assert.ErrorContains(t, err, "HTTP 502")
}

func TestClientGetCurrentRejectsOversizedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, strings.Repeat("x", maxResponseBodySize+1))
	}))
	defer server.Close()

	client, err := NewClient(Config{BaseURL: server.URL, Timeout: time.Second})
	require.NoError(t, err)

	_, err = client.GetCurrent(context.Background(), "Beijing")
	assert.ErrorContains(t, err, "exceeds size limit")
}

func TestClientGetCurrentPropagatesContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	client, err := NewClient(Config{BaseURL: server.URL, Timeout: time.Second})
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = client.GetCurrent(ctx, "Beijing")
	assert.True(t, errors.Is(err, context.Canceled))
}

func TestNewClientValidatesConfig(t *testing.T) {
	_, err := NewClient(Config{BaseURL: "not-a-url", Timeout: time.Second})
	assert.Error(t, err)

	_, err = NewClient(Config{BaseURL: "https://wttr.in"})
	assert.Error(t, err)
}
