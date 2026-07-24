package adapter

import "context"

// WeatherAdapter 定义天气业务依赖的外部能力边界。
// 你可以根据实际业务扩展方法参数和返回值
// 这里仅做演示
// Service 只依赖该接口，不感知外部服务协议和 HTTP 客户端实现。
type WeatherAdapter interface {
	GetWeather(ctx context.Context, city string) (string, error)
}
