package dao

import (
	"context"

	"go-api-template/internal/store/db"
	"go-api-template/internal/store/model"

	"gorm.io/gorm"
)

// TaskDAO 任务数据访问对象
type TaskDAO struct {
	db *gorm.DB
}

func NewTaskDAO() *TaskDAO {
	return &TaskDAO{db: db.GetGORM()}
}

func (d *TaskDAO) Create(ctx context.Context, task *model.Task) error {
	return d.db.WithContext(ctx).Create(task).Error
}

func (d *TaskDAO) List(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	err := d.db.WithContext(ctx).Order("id desc").Find(&tasks).Error
	return tasks, err
}

func (d *TaskDAO) GetByID(ctx context.Context, id uint) (*model.Task, error) {
	var task model.Task
	err := d.db.WithContext(ctx).First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (d *TaskDAO) Update(ctx context.Context, task *model.Task) error {
	return d.db.WithContext(ctx).Save(task).Error
}

func (d *TaskDAO) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&model.Task{}, id).Error
}
