// Code generated by "stringer -type ErrCode -linecomment"; DO NOT EDIT.

package error

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NoError-0]
	_ = x[GeneralError-199999]
	_ = x[ParamInvalid-101001]
	_ = x[ParamMissing-101002]
}

const (
	_ErrCode_name_0 = "no error"
	_ErrCode_name_1 = "invalid request parametermissing request parameter"
	_ErrCode_name_2 = "general error"
)

var (
	_ErrCode_index_1 = [...]uint8{0, 25, 50}
)

func (i ErrCode) String() string {
	switch {
	case i == 0:
		return _ErrCode_name_0
	case 101001 <= i && i <= 101002:
		i -= 101001
		return _ErrCode_name_1[_ErrCode_index_1[i]:_ErrCode_index_1[i+1]]
	case i == 199999:
		return _ErrCode_name_2
	default:
		return "ErrCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
