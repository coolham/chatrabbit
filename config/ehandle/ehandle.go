package ehandle

import "github.com/pkg/errors"

//go:generate stringer -type ErrCode -linecomment

// 1、自定义error结构体，并重写Error()方法

type ErrCode int64 //错误码

// 错误时返回自定义结构
type CustomError struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

func (e *CustomError) Error() string {
	return e.Code.String()
}

// 2、定义errorCode
const (
	CodeOk ErrCode = 0 // OK
)

// 3、新建自定义error实例化
func NewCustomError(code ErrCode) error {
	// 初次调用得用Wrap方法，进行实例化
	return errors.Wrap(&CustomError{
		Code:    code,
		Message: code.String(),
	}, "")
}
