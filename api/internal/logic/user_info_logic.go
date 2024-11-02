package logic

import (
	"context"
	"encoding/json"

	"api/internal/grpc_client/user"
	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (resp *types.UserInfoResponse, err error) {
	userId, err := l.ctx.Value(types.UserIdKey).(json.Number).Int64()
	if err != nil {
		l.Logger.Errorf("Get ctx userId error: err is %s", err.Error())
		return nil, err
	}
	if userId == 0 {
		return &types.UserInfoResponse{}, nil
	}
	res, err := l.svcCtx.UserRpc.FindById(l.ctx, &user.FindByIdRequest{
		UserId: userId,
	})
	if err != nil {
		l.Logger.Errorf("Find by userId rpc userId error: err is %s", err.Error())
		return nil, err
	}
	return &types.UserInfoResponse{
		UserId:   res.UserId,
		Username: res.Username,
		Avatar:   res.Avatar,
	}, nil
}
