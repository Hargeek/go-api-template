package config

// SQLiteConfig SQLite 数据库配置
type SQLiteConfig struct {
	Path    string `mapstructure:"path"`          // 数据库文件路径，如 ./data/app.db
	LogMode bool   `mapstructure:"log_mode"`      // 是否开启慢查询日志
	MaxIdle int    `mapstructure:"max_idle_conn"` // 最大空闲连接数（建议 1）
	MaxOpen int    `mapstructure:"max_open_conn"` // 最大连接数（建议 1）
	MaxLife int    `mapstructure:"max_life_time"` // 连接最大存活时间（秒）
}

// PostgresConfig PostgreSQL 数据库配置（切换时在 types.go 中取消注释）
type PostgresConfig struct {
	Host     string `mapstructure:"host"`          // 数据库地址
	Port     int64  `mapstructure:"port"`          // 数据库端口
	Database string `mapstructure:"db_name"`       // 数据库名
	Username string `mapstructure:"db_user"`       // 用户名
	Password string `mapstructure:"db_password"`   // 密码
	LogMode  bool   `mapstructure:"log_mode"`      // 是否开启慢查询日志
	MaxIdle  int    `mapstructure:"max_idle_conn"` // 最大空闲连接数
	MaxOpen  int    `mapstructure:"max_open_conn"` // 最大连接数
	MaxLife  int    `mapstructure:"max_life_time"` // 连接最大存活时间（秒）
}
