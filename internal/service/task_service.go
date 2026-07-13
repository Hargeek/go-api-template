package service

import (
	"context"

	errort "go-api-template/common/error"
	"go-api-template/internal/store/model"
)

// TaskService 任务业务逻辑接口
type TaskService interface {
	// Create 创建任务，返回创建后的任务（含自增 ID 和时间戳）
	Create(ctx context.Context, title, description string) (*model.Task, *errort.ApiError)
	// List 获取全部任务，按 ID 倒序排列
	List(ctx context.Context) ([]model.Task, *errort.ApiError)
	// GetByID 按 ID 获取任务详情，任务不存在时返回 TaskNotFound 错误
	GetByID(ctx context.Context, id uint) (*model.Task, *errort.ApiError)
	// Update 全量更新指定 ID 的任务，任务不存在时返回 TaskNotFound 错误
	Update(ctx context.Context, id uint, title, description string, done bool) (*model.Task, *errort.ApiError)
	// Delete 软删除指定 ID 的任务，任务不存在时不报错
	Delete(ctx context.Context, id uint) *errort.ApiError
}
