package types

// CreateTaskRequest 创建任务请求参数
type CreateTaskRequest struct {
	Title       string `json:"title"       binding:"required"`
	Description string `json:"description"`
}

// UpdateTaskRequest 更新任务请求参数
type UpdateTaskRequest struct {
	Title       string `json:"title"       binding:"required"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}
