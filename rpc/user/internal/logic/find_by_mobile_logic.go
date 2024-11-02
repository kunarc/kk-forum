package logic

import (
	"context"

	"user/internal/svc"
	"user/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByMobileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByMobileLogic {
	return &FindByMobileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByMobileLogic) FindByMobile(in *pb.FindByMobileRequest) (*pb.FindByMobileResponse, error) {
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)
	if err != nil {
		l.Logger.Errorf("FindByMobile mobile: %s error: %v", in.Mobile, err)
		return nil, err
	}
	if user == nil {
		return &pb.FindByMobileResponse{}, nil
	}
	return &pb.FindByMobileResponse{
		UserId:   int64(user.Id),
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}
