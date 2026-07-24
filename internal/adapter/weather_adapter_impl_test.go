package adapter

import (
	"context"
	"errors"
	"testing"

	"go-api-template/pkg/weather"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeWeatherClient struct {
	current *weather.CurrentWeather
	err     error
	city    string
}

func (f *fakeWeatherClient) GetCurrent(_ context.Context, city string) (*weather.CurrentWeather, error) {
	f.city = city
	return f.current, f.err
}

func TestWeatherAdapterImplGetWeather(t *testing.T) {
	client := &fakeWeatherClient{current: &weather.CurrentWeather{Description: "晴", TemperatureC: 25}}
	weatherAdapter := newWeatherAdapter(client)

	result, err := weatherAdapter.GetWeather(context.Background(), "Beijing")

	require.NoError(t, err)
	assert.Equal(t, "Beijing: 晴，25°C", result)
	assert.Equal(t, "Beijing", client.city)
}

func TestWeatherAdapterImplReturnsClientError(t *testing.T) {
	want := errors.New("weather client unavailable")
	weatherAdapter := newWeatherAdapter(&fakeWeatherClient{err: want})

	result, err := weatherAdapter.GetWeather(context.Background(), "Beijing")

	assert.Empty(t, result)
	assert.ErrorIs(t, err, want)
}
