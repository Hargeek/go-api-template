package controller

import (
	errort "go-api-template/common/error"
	res "go-api-template/common/types/response"
	"go-api-template/internal/service"
	"net/http"
	"strconv"

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
	tasks, err := t.Service.List()
	if err != nil {
		res.ApiResponse(c, http.StatusInternalServerError, errort.GeneralError, "获取任务列表失败", nil)
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
//	@Param			body	body		createTaskRequest	true	"任务信息"
//	@Success		200		{object}	res.CommonApiResponseData
//	@Router			/api/v1/tasks [post]
func (t *TaskController) CreateTask(c *gin.Context) {
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ApiResponse(c, http.StatusBadRequest, errort.ParamInvalid, err.Error(), nil)
		return
	}
	task, err := t.Service.Create(req.Title, req.Description)
	if err != nil {
		res.ApiResponse(c, http.StatusInternalServerError, errort.GeneralError, "创建任务失败", nil)
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
	task, err := t.Service.GetByID(id)
	if err != nil {
		res.ApiResponse(c, http.StatusNotFound, errort.GeneralError, "任务不存在", nil)
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
//	@Param			id		path		int					true	"任务 ID"
//	@Param			body	body		updateTaskRequest	true	"任务信息"
//	@Success		200		{object}	res.CommonApiResponseData
//	@Router			/api/v1/tasks/{id} [put]
func (t *TaskController) UpdateTask(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		res.ApiResponse(c, http.StatusBadRequest, errort.ParamInvalid, "无效的任务 ID", nil)
		return
	}
	var req updateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ApiResponse(c, http.StatusBadRequest, errort.ParamInvalid, err.Error(), nil)
		return
	}
	task, err := t.Service.Update(id, req.Title, req.Description, req.Done)
	if err != nil {
		res.ApiResponse(c, http.StatusInternalServerError, errort.GeneralError, "更新任务失败", nil)
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
	if err := t.Service.Delete(id); err != nil {
		res.ApiResponse(c, http.StatusInternalServerError, errort.GeneralError, "删除任务失败", nil)
		return
	}
	res.ApiResponse(c, http.StatusOK, errort.NoError, "删除成功", nil)
}

// --- 请求体结构 ---

type createTaskRequest struct {
	Title       string `json:"title"       binding:"required"`
	Description string `json:"description"`
}

type updateTaskRequest struct {
	Title       string `json:"title"       binding:"required"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// --- 工具函数 ---

func parseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
