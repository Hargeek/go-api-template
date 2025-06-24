package types

import (
	"github.com/gin-gonic/gin"
	errort "go-api-template/common/error"
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
