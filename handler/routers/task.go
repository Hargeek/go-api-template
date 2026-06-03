package routers

import (
	"go-api-template/handler/controller"

	"github.com/gin-gonic/gin"
)

func (r *router) InitTaskRouter(api *gin.RouterGroup) {
	tasks := api.Group("/tasks")
	{
		tasks.GET("", controller.Task.ListTasks)
		tasks.POST("", controller.Task.CreateTask)
		tasks.GET("/:id", controller.Task.GetTask)
		tasks.PUT("/:id", controller.Task.UpdateTask)
		tasks.DELETE("/:id", controller.Task.DeleteTask)
	}
}
