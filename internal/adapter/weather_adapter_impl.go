package adapter

import (
	"context"
	"fmt"

	"go-api-template/pkg/weather"
)

// weatherClient 是 WeatherAdapterImpl 使用的最小客户端能力，便于在 Adapter 测试中替换真实 HTTP 客户端。
type weatherClient interface {
	GetCurrent(ctx context.Context, city string) (*weather.CurrentWeather, error)
}

// WeatherAdapterImpl 使用天气客户端实现 WeatherAdapter，并负责将外部结果转换为项目内部语义。
type WeatherAdapterImpl struct {
	client weatherClient
}

var _ WeatherAdapter = (*WeatherAdapterImpl)(nil)
var _ weatherClient = (*weather.Client)(nil)

// NewWeatherAdapter 创建天气外部服务适配器。
func NewWeatherAdapter(client *weather.Client) *WeatherAdapterImpl {
	return newWeatherAdapter(client)
}

// newWeatherAdapter 接收最小客户端接口，仅用于隔离构造和支持单元测试注入。
func newWeatherAdapter(client weatherClient) *WeatherAdapterImpl {
	return &WeatherAdapterImpl{client: client}
}

// GetWeather 将 pkg/weather 返回的协议无关结果格式化为 WeatherService 当前需要的天气描述。
func (a *WeatherAdapterImpl) GetWeather(ctx context.Context, city string) (string, error) {
	current, err := a.client.GetCurrent(ctx, city)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s: %s，%g°C", city, current.Description, current.TemperatureC), nil
}
