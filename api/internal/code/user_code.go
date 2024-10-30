package code

import "pkg/xerror"

var (
	RegisterMobileEmpty        = xerror.New(10001, "注册手机号不能为空")
	VerificationCodeEmpty      = xerror.New(100002, "验证码不能为空")
	MobileHasRegistered        = xerror.New(100003, "手机号已经注册")
	LoginMobileEmpty           = xerror.New(100003, "手机号不能为空")
	RegisterPasswdEmpty        = xerror.New(100004, "密码不能为空")
	SendVerificationCodeExceed = xerror.New(100005, "请勿频繁发送验证码")
	MobileIsError              = xerror.New(100006, "请输入正确的手机号")
	SendVerCodeError           = xerror.New(100007, "发送验证码失败")
	VerificationCodeExpire     = xerror.New(100008, "验证码过期")
	VerificationCodeError      = xerror.New(100009, "验证码错误")
)

var (
	RegisterUserNameEmpty = xerror.New(200001, "注册用户名不能为空")
	RegisterUserNameSpace = xerror.New(200002, "注册用户名首尾不能含空格")
)
