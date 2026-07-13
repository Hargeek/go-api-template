package adapter

import "context"

// WeatherAdapter 定义外部天气服务的适配接口
// 你可以根据实际业务扩展方法参数和返回值
// 这里仅做演示

type WeatherAdapter interface {
	GetWeather(ctx context.Context, city string) (string, error)
}
