package logic

import (
	"context"
	"pkg/xerror"
	"time"

	"user/internal/code"
	"user/internal/model"
	"user/internal/svc"
	"user/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)
	if user != nil {
		return nil, code.RegisterMobileRepeat
	}
	if err != nil && err != sqlc.ErrNotFound {
		l.Logger.Errorf("find user by mobile error: err is %s", err.Error())
		return nil, xerror.ServerErr
	}
	res, err := l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Username:   in.Username,
		Mobile:     in.Mobile,
		Avatar:     in.Avatar,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("save user error: err is %s", err.Error())
		return nil, xerror.ServerErr
	}
	userId, _ := res.LastInsertId()
	return &pb.RegisterResponse{
		UserId: userId,
	}, nil
}
