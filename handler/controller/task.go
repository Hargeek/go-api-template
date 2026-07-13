package controller

import (
	"net/http"
	"strconv"

	errort "go-api-template/common/error"
	typeT "go-api-template/common/types"
	res "go-api-template/common/types/response"
	"go-api-template/internal/service"

	"github.com/gin-gonic/gin"
)

var Task *TaskController

type TaskController struct {
	Service service.TaskService
}

func NewTaskController(s service.TaskService) *TaskController {
	return &TaskController{Service: s}
}

// ListTasks 获取任务列表
//
//	@Accept			json
//	@Produce		json
//	@Summary		获取任务列表
//	@Description	获取所有任务，按创建时间倒序
//	@Tags			Task API
//	@Success		200	{object}	res.CommonApiResponseData
//	@Router			/api/v1/tasks [get]
func (t *TaskController) ListTasks(c *gin.Context) {
	tasks, apiErr := t.Service.List(c.Request.Context())
	if apiErr != nil {
		res.ApiResponse(c, http.StatusInternalServerError, apiErr.Code, apiErr.Msg, nil)
		return
	}
	res.ApiResponse(c, http.StatusOK, errort.NoError, "ok", tasks)
}

// CreateTask 创建任务
//
//	@Accept			json
//	@Produce		json
//	@Summary		创建任务
//	@Description	创建一个新任务
//	@Tags			Task API
//	@Param			body	body		typeT.CreateTaskRequest	true	"任务信息"
//	@Success		200		{object}	res.CommonApiResponseData
//	@Router			/api/v1/tasks [post]
func (t *TaskController) CreateTask(c *gin.Context) {
	var req typeT.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ApiResponse(c, http.StatusBadRequest, errort.ParamInvalid, err.Error(), nil)
		return
	}
	task, apiErr := t.Service.Create(c.Request.Context(), req.Title, req.Description)
	if apiErr != nil {
		res.ApiResponse(c, http.StatusInternalServerError, apiErr.Code, apiErr.Msg, nil)
		return
	}
	res.ApiResponse(c, http.StatusOK, errort.NoError, "创建成功", task)
}

// GetTask 获取单个任务
//
//	@Accept			json
//	@Produce		json
//	@Summary		获取任务详情
//	@Description	根据 ID 获取任务详情
//	@Tags			Task API
//	@Param			id	path		int	true	"任务 ID"
//	@Success		200	{object}	res.CommonApiResponseData
//	@Router			/api/v1/tasks/{id} [get]
func (t *TaskController) GetTask(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		res.ApiResponse(c, http.StatusBadRequest, errort.ParamInvalid, "无效的任务 ID", nil)
		return
	}
	task, apiErr := t.Service.GetByID(c.Request.Context(), id)
	if apiErr != nil {
		if apiErr.Code == errort.TaskNotFound {
			res.ApiResponse(c, http.StatusNotFound, apiErr.Code, apiErr.Msg, nil)
			return
		}
		res.ApiResponse(c, http.StatusInternalServerError, apiErr.Code, apiErr.Msg, nil)
		return
	}
	res.ApiResponse(c, http.StatusOK, errort.NoError, "ok", task)
}

// UpdateTask 更新任务
//
//	@Accept			json
//	@Produce		json
//	@Summary		更新任务
//	@Description	更新指定 ID 的任务
//	@Tags			Task API
//	@Param			id		path		int						true	"任务 ID"
//	@Param			body	body		typeT.UpdateTaskRequest	true	"任务信息"
//	@Success		200		{object}	res.CommonApiResponseData
//	@Router			/api/v1/tasks/{id} [put]
func (t *TaskController) UpdateTask(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		res.ApiResponse(c, http.StatusBadRequest, errort.ParamInvalid, "无效的任务 ID", nil)
		return
	}
	var req typeT.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ApiResponse(c, http.StatusBadRequest, errort.ParamInvalid, err.Error(), nil)
		return
	}
	task, apiErr := t.Service.Update(c.Request.Context(), id, req.Title, req.Description, req.Done)
	if apiErr != nil {
		if apiErr.Code == errort.TaskNotFound {
			res.ApiResponse(c, http.StatusNotFound, apiErr.Code, apiErr.Msg, nil)
			return
		}
		res.ApiResponse(c, http.StatusInternalServerError, apiErr.Code, apiErr.Msg, nil)
		return
	}
	res.ApiResponse(c, http.StatusOK, errort.NoError, "更新成功", task)
}

// DeleteTask 删除任务
//
//	@Accept			json
//	@Produce		json
//	@Summary		删除任务
//	@Description	软删除指定 ID 的任务
//	@Tags			Task API
//	@Param			id	path		int	true	"任务 ID"
//	@Success		200	{object}	res.CommonApiResponseData
//	@Router			/api/v1/tasks/{id} [delete]
func (t *TaskController) DeleteTask(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		res.ApiResponse(c, http.StatusBadRequest, errort.ParamInvalid, "无效的任务 ID", nil)
		return
	}
	if apiErr := t.Service.Delete(c.Request.Context(), id); apiErr != nil {
		res.ApiResponse(c, http.StatusInternalServerError, apiErr.Code, apiErr.Msg, nil)
		return
	}
	res.ApiResponse(c, http.StatusOK, errort.NoError, "删除成功", nil)
}

// --- 工具函数 ---

func parseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
