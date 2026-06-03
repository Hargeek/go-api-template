package model

import "gorm.io/gorm"

// Task 任务模型（API + DB 示例）
type Task struct {
	gorm.Model
	Title       string `gorm:"not null"        json:"title"`
	Description string `gorm:"default:''"      json:"description"`
	Done        bool   `gorm:"default:false"   json:"done"`
}
