package dao

import (
	"context"
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"go-api-template/internal/store/model"
)

func TestTaskDAOConditions(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.Task{}); err != nil {
		t.Fatal(err)
	}
	d := &Dao{db: db}
	ctx := context.Background()

	first := &model.Task{Title: "first"}
	second := &model.Task{Title: "second"}
	if err := d.CreateTask(ctx, first); err != nil {
		t.Fatal(err)
	}
	if err := d.CreateTask(ctx, second); err != nil {
		t.Fatal(err)
	}

	tasks, err := d.ListTasks(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 2 || tasks[0].ID != second.ID || tasks[1].ID != first.ID {
		t.Fatalf("ListTasks() order = %+v", tasks)
	}

	task, err := d.GetTaskByID(ctx, first.ID)
	if err != nil {
		t.Fatal(err)
	}
	if task.Title != first.Title {
		t.Fatalf("GetTaskByID() = %+v", task)
	}

	if err := d.DeleteTask(ctx, first.ID); err != nil {
		t.Fatal(err)
	}
	_, err = d.GetTaskByID(ctx, first.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("GetTaskByID() after delete error = %v", err)
	}
}
