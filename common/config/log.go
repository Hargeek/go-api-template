package config

type LogConfig struct {
	Level  string   `mapstructure:"level"`  // 日志级别
	Output []string `mapstructure:"output"` // 日志输出
}
