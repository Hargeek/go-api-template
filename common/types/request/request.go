package types

type CommonListRequestParams struct {
	Page  int `form:"page" json:"page" binding:"omitempty" example:"1"`    // 页码
	Limit int `form:"limit" json:"limit" binding:"omitempty" example:"10"` // 每页数量
}
