package service

import (
	"go-api-template/internal/store/dao"
	"go-api-template/internal/store/model"
)

type TaskServiceImpl struct {
	dao *dao.TaskDAO
}

func NewTaskServiceImpl() *TaskServiceImpl {
	return &TaskServiceImpl{dao: dao.NewTaskDAO()}
}

func (s *TaskServiceImpl) Create(title, description string) (*model.Task, error) {
	task := &model.Task{
		Title:       title,
		Description: description,
	}
	if err := s.dao.Create(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskServiceImpl) List() ([]model.Task, error) {
	return s.dao.List()
}

func (s *TaskServiceImpl) GetByID(id uint) (*model.Task, error) {
	return s.dao.GetByID(id)
}

func (s *TaskServiceImpl) Update(id uint, title, description string, done bool) (*model.Task, error) {
	task, err := s.dao.GetByID(id)
	if err != nil {
		return nil, err
	}
	task.Title = title
	task.Description = description
	task.Done = done
	if err := s.dao.Update(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskServiceImpl) Delete(id uint) error {
	return s.dao.Delete(id)
}
