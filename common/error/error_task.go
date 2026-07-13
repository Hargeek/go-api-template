package error

// 任务模块（102）错误信息定义，service 层用 fmt.Errorf 填充动态参数
const (
	MsgTaskNotFound     = "任务 %d 不存在"
	MsgTaskCreateFailed = "创建任务失败: %v"
	MsgTaskListFailed   = "获取任务列表失败: %v"
	MsgTaskGetFailed    = "获取任务 %d 失败: %v"
	MsgTaskUpdateFailed = "更新任务 %d 失败: %v"
	MsgTaskDeleteFailed = "删除任务 %d 失败: %v"
)
