package controller

import (
	errort "go-api-template/common/error"
	res "go-api-template/common/types/response"
	"go-api-template/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Hello *HelloController

type HelloController struct {
	Service service.HelloService
}

func NewHelloController(s service.HelloService) *HelloController {
	return &HelloController{Service: s}
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
func (h *HelloController) HelloController(c *gin.Context) {
	res.ApiResponse(c, http.StatusOK, errort.NoError, h.Service.Hello(), nil)
}
