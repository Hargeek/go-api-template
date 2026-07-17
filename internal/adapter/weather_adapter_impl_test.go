package adapter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWeatherAdapterImpl_GetWeather(t *testing.T) {
	weatherAdapter := NewWeatherAdapterImpl("test-api-key")

	result, err := weatherAdapter.GetWeather(context.Background(), "Beijing")

	require.NoError(t, err)
	assert.Equal(t, "test-api-key", weatherAdapter.ApiKey)
	assert.Equal(t, "Beijing: 晴，25°C", result)
}
