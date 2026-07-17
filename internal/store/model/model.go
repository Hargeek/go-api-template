package model

import "time"

// Model 定义项目通用的持久化基础字段。
type Model struct {
	ID        uint       `json:"id" gorm:"primaryKey;index"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}
