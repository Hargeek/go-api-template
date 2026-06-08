package service

import (
	"context"

	"go-api-template/internal/store/model"
)

// TaskService 任务业务逻辑接口
type TaskService interface {
	Create(ctx context.Context, title, description string) (*model.Task, error)
	List(ctx context.Context) ([]model.Task, error)
	GetByID(ctx context.Context, id uint) (*model.Task, error)
	Update(ctx context.Context, id uint, title, description string, done bool) (*model.Task, error)
	Delete(ctx context.Context, id uint) error
}
