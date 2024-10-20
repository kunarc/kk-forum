package xcode

import "strconv"

// type
//
//	@param error interface {
//	@author kunarc
//	@update 2024-10-19 05:49:38
type (
	XError interface {
		Error() string
		Code() int
		Message() string
	}
	XCode struct {
		code int
		msg  string
	}
)

func New(code int, msg string) XError {
	return &XCode{
		code: code,
		msg:  msg,
	}
}

// Error
//
//	@receiver x *XCode
//	@return string
//	@author kunarc
//	@update 2024-10-19 05:53:09
func (x *XCode) Error() string {
	if len(x.msg) > 0 {
		return x.msg
	}
	return strconv.Itoa(x.code)
}

// Code
//
//	@receiver x *XCode
//	@return int
//	@author kunarc
//	@update 2024-10-19 05:53:44
func (x *XCode) Code() int {
	return x.code
}

// Message
//
//	@receiver x *XCode
//	@return string
//	@author kunarc
//	@update 2024-10-19 05:54:17
func (x *XCode) Message() string {
	return x.Error()
}