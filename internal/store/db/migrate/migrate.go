package migrate

import (
	"go-api-template/internal/store/db"
	"go-api-template/internal/store/model"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() {
	if err := db.GetGORM().AutoMigrate(
		&model.Task{},
	); err != nil {
		panic("auto migrate failed: " + err.Error())
	}
}
