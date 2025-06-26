package service

import "go-api-template/internal/adapter"

type WeatherServiceImpl struct {
	Adapter adapter.WeatherAdapter
}

func NewWeatherServiceImpl(adapter adapter.WeatherAdapter) *WeatherServiceImpl {
	return &WeatherServiceImpl{Adapter: adapter}
}

func (s *WeatherServiceImpl) QueryWeather(city string) (string, error) {
	return s.Adapter.GetWeather(city)
}
