package dao

import (
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

func (d *TaskDAO) Create(task *model.Task) error {
	return d.db.Create(task).Error
}

func (d *TaskDAO) List() ([]model.Task, error) {
	var tasks []model.Task
	err := d.db.Order("id desc").Find(&tasks).Error
	return tasks, err
}

func (d *TaskDAO) GetByID(id uint) (*model.Task, error) {
	var task model.Task
	err := d.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (d *TaskDAO) Update(task *model.Task) error {
	return d.db.Save(task).Error
}

func (d *TaskDAO) Delete(id uint) error {
	return d.db.Delete(&model.Task{}, id).Error
}
