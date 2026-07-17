package service

import (
	"context"
	"errors"
	"testing"

	errort "go-api-template/common/error"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockWeatherAdapter struct {
	mock.Mock
}

func (m *mockWeatherAdapter) GetWeather(ctx context.Context, city string) (string, error) {
	args := m.Called(ctx, city)
	return args.String(0), args.Error(1)
}

func TestWeatherServiceImpl_QueryWeather(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		weatherAdapter := new(mockWeatherAdapter)
		weatherAdapter.On("GetWeather", mock.Anything, "Beijing").
			Return("Beijing: 晴，25°C", nil).
			Once()
		weatherService := NewWeatherServiceImpl(weatherAdapter)

		result, apiErr := weatherService.QueryWeather(context.Background(), "Beijing")

		require.Nil(t, apiErr)
		assert.Equal(t, "Beijing: 晴，25°C", result)
		weatherAdapter.AssertExpectations(t)
	})

	t.Run("adapter error", func(t *testing.T) {
		weatherAdapter := new(mockWeatherAdapter)
		weatherAdapter.On("GetWeather", mock.Anything, "Unknown").
			Return("", errors.New("upstream unavailable")).
			Once()
		weatherService := NewWeatherServiceImpl(weatherAdapter)

		result, apiErr := weatherService.QueryWeather(context.Background(), "Unknown")

		assert.Empty(t, result)
		require.NotNil(t, apiErr)
		assert.Equal(t, errort.GeneralError, apiErr.Code)
		assert.Contains(t, apiErr.Msg, "upstream unavailable")
		weatherAdapter.AssertExpectations(t)
	})
}
