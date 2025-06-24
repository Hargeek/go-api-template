package controller

import (
	"github.com/gin-gonic/gin"
	errort "go-api-template/common/error"
	res "go-api-template/common/types/response"
	"go-api-template/internal/service"
	"net/http"
)

var Hello helloController

type helloController struct {
	Service *service.HelloServiceImpl
}

// HelloController
//
// @Accept      json
// @Produce     json
// @Summary     Hello World 接口
// @Description Hello World 接口
// @Tags        Hello API
// @Success     200 {object} res.CommonApiResponseData
// @Router      /api/v1/hello [get]
func (h *helloController) HelloController(c *gin.Context) {
	res.ApiResponse(c, http.StatusOK, errort.NoError, h.Service.Hello(), nil)
}
