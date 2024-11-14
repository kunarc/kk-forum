package code

import "pkg/xerror"

var (
	FollowUserIdEmpty         = xerror.New(40001, "关注用户id为空")
	FollowedUserIdEmpty       = xerror.New(40002, "被关注用户id为空")
	CannotFollowSelf          = xerror.New(40003, "不能关注自己")
	UserIdEmpty               = xerror.New(40004, "用户id为空")
	CannotCancelSelf          = xerror.New(40005, "用户不能取关自己")
	CancelUserIdEmpty         = xerror.New(40006, "取关用户为空")
	CancelFollowedUserIdEmpty = xerror.New(40007, "被取关用户为空")
	CancelObjectInvaild       = xerror.New(40008, "取关对象不合法")
)
