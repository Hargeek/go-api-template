package dao

import (
	"context"

	"go-api-template/internal/store/dao/base"
	"go-api-template/internal/store/model"
)

// CreateTask 插入一条任务记录，成功后 task 会回填自增 ID 和时间戳
func (d *Dao) CreateTask(ctx context.Context, task *model.Task) error {
	return d.db.WithContext(ctx).Create(task).Error
}

// ListTasks 查询全部任务，按 ID 倒序排列
func (d *Dao) ListTasks(ctx context.Context) ([]model.Task, error) {
	return d.taskList(ctx, &base.DBConditions{Order: "id DESC"})
}

// GetTaskByID 按 ID 查询单条任务，记录不存在时返回 gorm.ErrRecordNotFound
func (d *Dao) GetTaskByID(ctx context.Context, id uint) (*model.Task, error) {
	return d.getTask(ctx, &base.DBConditions{
		And: map[string]interface{}{"id = ?": id},
	})
}

// UpdateTask 按主键全量保存任务记录
func (d *Dao) UpdateTask(ctx context.Context, task *model.Task) error {
	return d.db.WithContext(ctx).Save(task).Error
}

// DeleteTask 按 ID 软删除任务（填充 deleted_at），记录不存在时不报错
func (d *Dao) DeleteTask(ctx context.Context, id uint) error {
	_, err := d.taskDelete(ctx, &base.DBConditions{
		And: map[string]interface{}{"id = ?": id},
	})
	return err
}

// taskList 执行任务列表查询，统一应用查询条件。
func (d *Dao) taskList(ctx context.Context, conditions *base.DBConditions) ([]model.Task, error) {
	tasks := make([]model.Task, 0)
	db := d.db.WithContext(ctx).Model(&model.Task{})
	db = conditions.Fill(db)
	if err := db.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// getTask 执行任务单条查询，记录不存在时保留 gorm.ErrRecordNotFound 语义。
func (d *Dao) getTask(ctx context.Context, conditions *base.DBConditions) (*model.Task, error) {
	task := &model.Task{}
	db := d.db.WithContext(ctx).Model(&model.Task{})
	db = conditions.Fill(db)
	if err := db.First(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

// taskDelete 按条件软删除任务并返回影响行数。
func (d *Dao) taskDelete(ctx context.Context, conditions *base.DBConditions) (int64, error) {
	db := d.db.WithContext(ctx).Model(&model.Task{})
	db = conditions.Fill(db)
	result := db.Delete(&model.Task{})
	return result.RowsAffected, result.Error
}
