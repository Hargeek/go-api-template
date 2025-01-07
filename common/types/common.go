package types

type ListParams struct {
	Page  int `form:"page" json:"page" binding:"omitempty" example:"1"`    // 页码
	Limit int `form:"limit" json:"limit" binding:"omitempty" example:"10"` // 每页数量
}

type CommonApiResponse struct {
	Msg  string      `json:"msg"`  // message
	Data interface{} `json:"data"` // data
	Code int         `json:"code"` // code
}
