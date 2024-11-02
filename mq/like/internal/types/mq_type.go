package types

type (
	ThumbupMqMsg struct {
		BizId    string // 业务id
		ObjId    int64  // 点赞对象id
		UserId   int64  // 用户id
		LikeType int32
	}
)
