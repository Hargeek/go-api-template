package service

import (
	"go-api-template/common/metrics"
	"go-api-template/internal/adapter"
)

type WeatherServiceImpl struct {
	Adapter adapter.WeatherAdapter
}

func NewWeatherServiceImpl(adapter adapter.WeatherAdapter) *WeatherServiceImpl {
	return &WeatherServiceImpl{Adapter: adapter}
}

func (s *WeatherServiceImpl) QueryWeather(city string) (string, error) {
	result, err := s.Adapter.GetWeather(city)
	if err != nil {
		metrics.WeatherQueryTotal.WithLabelValues(city, "fail").Inc()
		return "", err
	}
	metrics.WeatherQueryTotal.WithLabelValues(city, "success").Inc()
	return result, nil
}
