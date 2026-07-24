package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	errort "go-api-template/common/error"
	"go-api-template/common/logger"
	// profile:mtl:start
	"go-api-template/common/metrics"
	// profile:mtl:end
	"go-api-template/internal/store/dao"
	"go-api-template/internal/store/model"
)

type TaskServiceImpl struct {
	dao *dao.Dao
}

func NewTaskServiceImpl(d *dao.Dao) *TaskServiceImpl {
	return &TaskServiceImpl{dao: d}
}

func (s *TaskServiceImpl) Create(ctx context.Context, title, description string) (*model.Task, *errort.ApiError) {
	task := &model.Task{
		Title:       title,
		Description: description,
	}
	if err := s.dao.CreateTask(ctx, task); err != nil {
		// profile:mtl:start
		metrics.TaskOperationsTotal.WithLabelValues("create", "fail").Inc()
		// profile:mtl:end
		logger.ErrorContext(ctx, "task create failed", "error", err)
		return nil, errort.NewApiError(errort.GeneralError, fmt.Errorf(errort.MsgTaskCreateFailed, err))
	}
	// profile:mtl:start
	metrics.TaskOperationsTotal.WithLabelValues("create", "success").Inc()
	// profile:mtl:end
	return task, nil
}

func (s *TaskServiceImpl) List(ctx context.Context) ([]model.Task, *errort.ApiError) {
	tasks, err := s.dao.ListTasks(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "task list failed", "error", err)
		return nil, errort.NewApiError(errort.GeneralError, fmt.Errorf(errort.MsgTaskListFailed, err))
	}
	return tasks, nil
}

func (s *TaskServiceImpl) GetByID(ctx context.Context, id uint) (*model.Task, *errort.ApiError) {
	task, err := s.dao.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errort.NewApiError(errort.TaskNotFound, fmt.Errorf(errort.MsgTaskNotFound, id))
		}
		logger.ErrorContext(ctx, "task get failed",
			"task_id", id,
			"error", err,
		)
		return nil, errort.NewApiError(errort.GeneralError, fmt.Errorf(errort.MsgTaskGetFailed, id, err))
	}
	return task, nil
}

func (s *TaskServiceImpl) Update(ctx context.Context, id uint, title, description string, done bool) (*model.Task, *errort.ApiError) {
	task, apiErr := s.GetByID(ctx, id)
	if apiErr != nil {
		return nil, apiErr
	}
	task.Title = title
	task.Description = description
	task.Done = done
	if err := s.dao.UpdateTask(ctx, task); err != nil {
		logger.ErrorContext(ctx, "task update failed",
			"task_id", id,
			"error", err,
		)
		return nil, errort.NewApiError(errort.GeneralError, fmt.Errorf(errort.MsgTaskUpdateFailed, id, err))
	}
	return task, nil
}

func (s *TaskServiceImpl) Delete(ctx context.Context, id uint) *errort.ApiError {
	if err := s.dao.DeleteTask(ctx, id); err != nil {
		logger.ErrorContext(ctx, "task delete failed",
			"task_id", id,
			"error", err,
		)
		return errort.NewApiError(errort.GeneralError, fmt.Errorf(errort.MsgTaskDeleteFailed, id, err))
	}
	return nil
}
