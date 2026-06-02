package types

import (
	errort "go-api-template/common/error"

	"github.com/gin-gonic/gin"
)

type CommonApiResponseData struct {
	Msg  string         `json:"msg"`  // message
	Data interface{}    `json:"data"` // data
	Code errort.ErrCode `json:"code"` // code
}

func ApiResponse(c *gin.Context, httpCode int, apiCode errort.ErrCode, msg string, data interface{}) {
	c.JSON(httpCode, CommonApiResponseData{
		Msg:  msg,
		Data: data,
		Code: apiCode,
	})
}
