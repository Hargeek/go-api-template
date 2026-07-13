package adapter

import (
	"context"
	"fmt"
)

// WeatherAdapterImpl 是 WeatherAdapter 的一个简单实现
// 实际项目中可调用第三方API（用 http.NewRequestWithContext 透传 ctx），这里仅做演示

type WeatherAdapterImpl struct {
	ApiKey string
}

func NewWeatherAdapterImpl(apiKey string) *WeatherAdapterImpl {
	return &WeatherAdapterImpl{ApiKey: apiKey}
}

func (w *WeatherAdapterImpl) GetWeather(_ context.Context, city string) (string, error) {
	return fmt.Sprintf("%s: 晴，25°C", city), nil
}
