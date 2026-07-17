package controller

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	errort "go-api-template/common/error"
	"go-api-template/internal/store/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type stubTaskService struct {
	create func(context.Context, string, string) (*model.Task, *errort.ApiError)
	list   func(context.Context) ([]model.Task, *errort.ApiError)
	get    func(context.Context, uint) (*model.Task, *errort.ApiError)
	update func(context.Context, uint, string, string, bool) (*model.Task, *errort.ApiError)
	delete func(context.Context, uint) *errort.ApiError
}

func (s *stubTaskService) Create(ctx context.Context, title, description string) (*model.Task, *errort.ApiError) {
	return s.create(ctx, title, description)
}

func (s *stubTaskService) List(ctx context.Context) ([]model.Task, *errort.ApiError) {
	return s.list(ctx)
}

func (s *stubTaskService) GetByID(ctx context.Context, id uint) (*model.Task, *errort.ApiError) {
	return s.get(ctx, id)
}

func (s *stubTaskService) Update(ctx context.Context, id uint, title, description string, done bool) (*model.Task, *errort.ApiError) {
	return s.update(ctx, id, title, description, done)
}

func (s *stubTaskService) Delete(ctx context.Context, id uint) *errort.ApiError {
	return s.delete(ctx, id)
}

func TestTaskController_ListTasks(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := &stubTaskService{
			list: func(context.Context) ([]model.Task, *errort.ApiError) {
				return []model.Task{{Title: "first"}}, nil
			},
		}

		response := performTaskRequest(http.MethodGet, "/tasks", "", "/tasks", NewTaskController(service).ListTasks)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Contains(t, response.Body.String(), `"title":"first"`)
	})

	t.Run("service error", func(t *testing.T) {
		service := &stubTaskService{
			list: func(context.Context) ([]model.Task, *errort.ApiError) {
				return nil, errort.NewApiError(errort.GeneralError, nil)
			},
		}

		response := performTaskRequest(http.MethodGet, "/tasks", "", "/tasks", NewTaskController(service).ListTasks)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestTaskController_CreateTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := &stubTaskService{
			create: func(_ context.Context, title, description string) (*model.Task, *errort.ApiError) {
				return &model.Task{Title: title, Description: description}, nil
			},
		}

		response := performTaskRequest(
			http.MethodPost,
			"/tasks",
			`{"title":"write tests","description":"controller"}`,
			"/tasks",
			NewTaskController(service).CreateTask,
		)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Contains(t, response.Body.String(), `"title":"write tests"`)
	})

	t.Run("invalid body", func(t *testing.T) {
		service := &stubTaskService{}

		response := performTaskRequest(
			http.MethodPost,
			"/tasks",
			`{"description":"missing title"}`,
			"/tasks",
			NewTaskController(service).CreateTask,
		)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Contains(t, response.Body.String(), `"code":101001`)
	})

	t.Run("service error", func(t *testing.T) {
		service := &stubTaskService{
			create: func(context.Context, string, string) (*model.Task, *errort.ApiError) {
				return nil, errort.NewApiError(errort.GeneralError, nil)
			},
		}

		response := performTaskRequest(
			http.MethodPost,
			"/tasks",
			`{"title":"write tests"}`,
			"/tasks",
			NewTaskController(service).CreateTask,
		)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestTaskController_GetTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := &stubTaskService{
			get: func(_ context.Context, id uint) (*model.Task, *errort.ApiError) {
				return &model.Task{Title: "task"}, nil
			},
		}

		response := performTaskRequest(http.MethodGet, "/tasks/1", "", "/tasks/:id", NewTaskController(service).GetTask)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		response := performTaskRequest(http.MethodGet, "/tasks/invalid", "", "/tasks/:id", NewTaskController(&stubTaskService{}).GetTask)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("not found", func(t *testing.T) {
		service := &stubTaskService{
			get: func(context.Context, uint) (*model.Task, *errort.ApiError) {
				return nil, errort.NewApiError(errort.TaskNotFound, nil)
			},
		}

		response := performTaskRequest(http.MethodGet, "/tasks/999", "", "/tasks/:id", NewTaskController(service).GetTask)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("service error", func(t *testing.T) {
		service := &stubTaskService{
			get: func(context.Context, uint) (*model.Task, *errort.ApiError) {
				return nil, errort.NewApiError(errort.GeneralError, nil)
			},
		}

		response := performTaskRequest(http.MethodGet, "/tasks/1", "", "/tasks/:id", NewTaskController(service).GetTask)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestTaskController_UpdateTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := &stubTaskService{
			update: func(_ context.Context, id uint, title, description string, done bool) (*model.Task, *errort.ApiError) {
				return &model.Task{Title: title, Description: description, Done: done}, nil
			},
		}

		response := performTaskRequest(
			http.MethodPut,
			"/tasks/1",
			`{"title":"updated","description":"controller","done":true}`,
			"/tasks/:id",
			NewTaskController(service).UpdateTask,
		)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Contains(t, response.Body.String(), `"done":true`)
	})

	t.Run("invalid body", func(t *testing.T) {
		response := performTaskRequest(
			http.MethodPut,
			"/tasks/1",
			`{"done":true}`,
			"/tasks/:id",
			NewTaskController(&stubTaskService{}).UpdateTask,
		)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		response := performTaskRequest(
			http.MethodPut,
			"/tasks/invalid",
			`{"title":"updated"}`,
			"/tasks/:id",
			NewTaskController(&stubTaskService{}).UpdateTask,
		)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("not found", func(t *testing.T) {
		service := &stubTaskService{
			update: func(context.Context, uint, string, string, bool) (*model.Task, *errort.ApiError) {
				return nil, errort.NewApiError(errort.TaskNotFound, nil)
			},
		}

		response := performTaskRequest(
			http.MethodPut,
			"/tasks/999",
			`{"title":"updated"}`,
			"/tasks/:id",
			NewTaskController(service).UpdateTask,
		)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("service error", func(t *testing.T) {
		service := &stubTaskService{
			update: func(context.Context, uint, string, string, bool) (*model.Task, *errort.ApiError) {
				return nil, errort.NewApiError(errort.GeneralError, nil)
			},
		}

		response := performTaskRequest(
			http.MethodPut,
			"/tasks/1",
			`{"title":"updated"}`,
			"/tasks/:id",
			NewTaskController(service).UpdateTask,
		)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestTaskController_DeleteTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := &stubTaskService{
			delete: func(context.Context, uint) *errort.ApiError { return nil },
		}

		response := performTaskRequest(http.MethodDelete, "/tasks/1", "", "/tasks/:id", NewTaskController(service).DeleteTask)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("service error", func(t *testing.T) {
		service := &stubTaskService{
			delete: func(context.Context, uint) *errort.ApiError {
				return errort.NewApiError(errort.GeneralError, nil)
			},
		}

		response := performTaskRequest(http.MethodDelete, "/tasks/1", "", "/tasks/:id", NewTaskController(service).DeleteTask)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		response := performTaskRequest(http.MethodDelete, "/tasks/invalid", "", "/tasks/:id", NewTaskController(&stubTaskService{}).DeleteTask)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func performTaskRequest(method, target, body, route string, handler gin.HandlerFunc) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Handle(method, route, handler)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	router.ServeHTTP(recorder, request)
	return recorder
}
