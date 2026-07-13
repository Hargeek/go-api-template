package error

// 天气模块错误信息定义，service 层用 fmt.Errorf 填充动态参数
const (
	MsgWeatherQueryFailed = "查询天气失败: %v"
)
