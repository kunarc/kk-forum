package xerror

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

// ErrorHandler
//
//	@param err error
//	@return int
//	@return any
//	@author kunarc
//	@update 2024-10-19 06:13:52
func ErrorHandler(err error) (int, any) {
	xe := XErrorFromError(err)
	return http.StatusOK, &XCode{
		code: xe.Code(),
		msg:  xe.Message(),
	}
}

// XErrorFromError
//
//	@param err error
//	@return xe XError
//	@author kunarc
//	@update 2024-10-19 06:13:55
func XErrorFromError(err error) (xe XError) {
	err = errors.Cause(err)
	if v, ok := err.(XError); ok {
		xe = v
		return
	}
	// 判断是否是客户端超时或取消
	switch err {
	case context.Canceled:
		xe = Canceled
	case context.DeadlineExceeded:
		xe = Deadline
	default:
		xe = ServerErr
	}
	return
}
