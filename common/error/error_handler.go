package error

type ApiError struct {
	Code ErrCode `json:"code"` // 错误码
	Msg  string  `json:"msg"`  // 错误信息
}

func (e *ApiError) Error() string {
	return e.Code.String()
}

func NewApiError(code ErrCode, originError error) *ApiError {
	msg := ""
	if originError != nil {
		msg = originError.Error()
	} else {
		msg = code.String()
	}
	return &ApiError{
		Code: code,
		Msg:  msg,
	}
}
