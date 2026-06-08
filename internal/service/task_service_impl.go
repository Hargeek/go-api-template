package service

import (
	"context"

	"go-api-template/common/metrics"
	"go-api-template/internal/store/dao"
	"go-api-template/internal/store/model"
)

type TaskServiceImpl struct {
	dao *dao.TaskDAO
}

func NewTaskServiceImpl() *TaskServiceImpl {
	return &TaskServiceImpl{dao: dao.NewTaskDAO()}
}

func (s *TaskServiceImpl) Create(ctx context.Context, title, description string) (*model.Task, error) {
	task := &model.Task{Title: title, Description: description}
	if err := s.dao.Create(ctx, task); err != nil {
		metrics.TaskOperationsTotal.WithLabelValues("create", "fail").Inc()
		return nil, err
	}
	metrics.TaskOperationsTotal.WithLabelValues("create", "success").Inc()
	return task, nil
}

func (s *TaskServiceImpl) List(ctx context.Context) ([]model.Task, error) {
	return s.dao.List(ctx)
}

func (s *TaskServiceImpl) GetByID(ctx context.Context, id uint) (*model.Task, error) {
	return s.dao.GetByID(ctx, id)
}

func (s *TaskServiceImpl) Update(ctx context.Context, id uint, title, description string, done bool) (*model.Task, error) {
	task, err := s.dao.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	task.Title = title
	task.Description = description
	task.Done = done
	if err := s.dao.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskServiceImpl) Delete(ctx context.Context, id uint) error {
	return s.dao.Delete(ctx, id)
}
