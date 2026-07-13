package dao

import (
	"sync"

	"gorm.io/gorm"

	"go-api-template/internal/store/db"
)

var (
	dao  *Dao
	once sync.Once
)

// Dao 数据访问层统一入口，所有领域的 DB 操作均挂载在此 struct 上
type Dao struct {
	db *gorm.DB
}

// NewDao 返回全局唯一的 Dao 实例，仅首次调用时的 db 参数生效
func NewDao(db *gorm.DB) *Dao {
	once.Do(func() {
		dao = &Dao{db: db}
	})
	return dao
}

// DB 返回底层 gorm.DB 实例，供事务等需要直接操作 DB 的场景使用
func (d *Dao) DB() *gorm.DB {
	return d.db
}

// Close 关闭底层数据库连接，服务优雅退出时调用
func (d *Dao) Close() error {
	return db.Close()
}
