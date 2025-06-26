//go:build ignore

package error

import (
	"testing"
)

func TestErrCode_String(t *testing.T) {
	tests := []struct {
		code ErrCode
		want string
	}{
		{NoError, "无错误"},
		{ParamInvalid, "无效参数"},
		{ParamMissing, "缺少参数"},
		{GeneralError, "通用错误"},
		{ErrCode(99999), "ErrCode(99999)"},
	}

	for _, tt := range tests {
		t.Run(tt.code.String(), func(t *testing.T) {
			if got := tt.code.String(); got != tt.want {
				t.Errorf("ErrCode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
