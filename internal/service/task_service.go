package service

import "go-api-template/internal/store/model"

// TaskService 任务业务逻辑接口
type TaskService interface {
	Create(title, description string) (*model.Task, error)
	List() ([]model.Task, error)
	GetByID(id uint) (*model.Task, error)
	Update(id uint, title, description string, done bool) (*model.Task, error)
	Delete(id uint) error
}
