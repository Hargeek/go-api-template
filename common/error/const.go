package error

type ErrCode int

// 自定义响应体错误码
/*
自定义错误码共六位，前三位表示模块，后三位表示错误，规则声明如下
* 0: 无错误；
* 1xx: 系统通用错误，不与业务直接绑定
  * 101: 参数错误如query或param类型错误

1、安装 stringer 工具
go install golang.org/x/tools/cmd/stringer
2、生成错误码字符串
make generate-error
*/

//go:generate stringer -type ErrCode -linecomment
const (
	NoError      ErrCode = 0      // no error
	GeneralError ErrCode = 199999 // general error
)

//go:generate stringer -type ErrCode -linecomment
const (
	ParamInvalid ErrCode = 101001 // invalid request parameter
	ParamMissing ErrCode = 101002 // missing request parameter
)
