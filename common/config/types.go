package config

type Config struct {
	EnvConfig    string       `mapstructure:"env"`    // 环境标识
	ServerConfig ServerConfig `mapstructure:"server"` // 服务配置
	SQLiteConfig SQLiteConfig `mapstructure:"sqlite"` // 数据库配置（当前使用 SQLite）
	// PostgresConfig PostgresConfig `mapstructure:"postgres"` // 切换 PostgreSQL
	LogConfig     LogConfig     `mapstructure:"logging"` // 日志配置
	WeatherConfig WeatherConfig `mapstructure:"weather"` // 天气示例的外部服务配置
}
