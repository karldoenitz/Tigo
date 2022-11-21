package binding

import (
	"fmt"
)

// NewTagErr 提供了一个错误类型闭包
func NewTagErr(formatStr string) func(string) error {
	return func(fieldName string) error {
		return fmt.Errorf(formatStr, fieldName)
	}
}

// RequiredErr 等变量表示不同的错误类型
var (
	RequiredErr = NewTagErr("%s is required")
	DefaultErr  = NewTagErr("%s default is invalid")
	RegexErr    = NewTagErr("%s regex can not match")
	InvalidErr  = NewTagErr("%s is invalid")
	UnDefineErr = NewTagErr("%s is undefined")
)
