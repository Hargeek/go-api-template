package dao

import (
	"context"

	"go-api-template/internal/store/model"
)

// CreateTask 插入一条任务记录，成功后 task 会回填自增 ID 和时间戳
func (d *Dao) CreateTask(ctx context.Context, task *model.Task) error {
	return d.db.WithContext(ctx).Create(task).Error
}

// ListTasks 查询全部任务，按 ID 倒序排列
func (d *Dao) ListTasks(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	err := d.db.WithContext(ctx).Order("id desc").Find(&tasks).Error
	return tasks, err
}

// GetTaskByID 按 ID 查询单条任务，记录不存在时返回 gorm.ErrRecordNotFound
func (d *Dao) GetTaskByID(ctx context.Context, id uint) (*model.Task, error) {
	var task model.Task
	err := d.db.WithContext(ctx).First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateTask 按主键全量保存任务记录
func (d *Dao) UpdateTask(ctx context.Context, task *model.Task) error {
	return d.db.WithContext(ctx).Save(task).Error
}

// DeleteTask 按 ID 软删除任务（填充 deleted_at），记录不存在时不报错
func (d *Dao) DeleteTask(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&model.Task{}, id).Error
}
