package service

import (
	"context"
	"testing"

	errort "go-api-template/common/error"
	"go-api-template/internal/store/dao"
	"go-api-template/internal/store/model"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestTaskServiceImpl(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, sqlDB.Close()) })
	require.NoError(t, db.AutoMigrate(&model.Task{}))
	taskService := NewTaskServiceImpl(dao.NewDao(db))
	ctx := context.Background()

	t.Run("create", func(t *testing.T) {
		resetTaskStore(t, db)
		task, apiErr := taskService.Create(ctx, "write tests", "cover service behavior")

		require.Nil(t, apiErr)
		require.NotNil(t, task)
		assert.NotZero(t, task.ID)
		assert.Equal(t, "write tests", task.Title)
		assert.Equal(t, "cover service behavior", task.Description)
	})

	t.Run("list", func(t *testing.T) {
		resetTaskStore(t, db)
		require.NoError(t, db.Create(&model.Task{Title: "listed task"}).Error)
		tasks, apiErr := taskService.List(ctx)

		require.Nil(t, apiErr)
		require.Len(t, tasks, 1)
		assert.Equal(t, "listed task", tasks[0].Title)
	})

	t.Run("get not found", func(t *testing.T) {
		resetTaskStore(t, db)
		task, apiErr := taskService.GetByID(ctx, 999)

		assert.Nil(t, task)
		require.NotNil(t, apiErr)
		assert.Equal(t, errort.TaskNotFound, apiErr.Code)
	})

	t.Run("update", func(t *testing.T) {
		resetTaskStore(t, db)
		taskToUpdate := &model.Task{Title: "before update"}
		require.NoError(t, db.Create(taskToUpdate).Error)

		task, apiErr := taskService.Update(ctx, taskToUpdate.ID, "tests completed", "updated", true)

		require.Nil(t, apiErr)
		require.NotNil(t, task)
		assert.Equal(t, "tests completed", task.Title)
		assert.Equal(t, "updated", task.Description)
		assert.True(t, task.Done)
	})

	t.Run("delete", func(t *testing.T) {
		resetTaskStore(t, db)
		taskToDelete := &model.Task{Title: "delete me"}
		require.NoError(t, db.Create(taskToDelete).Error)

		apiErr := taskService.Delete(ctx, taskToDelete.ID)
		require.Nil(t, apiErr)

		task, apiErr := taskService.GetByID(ctx, taskToDelete.ID)
		assert.Nil(t, task)
		require.NotNil(t, apiErr)
		assert.Equal(t, errort.TaskNotFound, apiErr.Code)
	})
}

func resetTaskStore(t *testing.T, db *gorm.DB) {
	t.Helper()
	require.NoError(t, db.Unscoped().Where("id > ?", 0).Delete(&model.Task{}).Error)
}
