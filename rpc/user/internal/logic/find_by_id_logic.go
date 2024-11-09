package logic

import (
	"context"

	"user/internal/svc"
	"user/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByIdLogic {
	return &FindByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByIdLogic) FindById(in *pb.FindByIdRequest) (*pb.FindByIdResponse, error) {
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		l.Logger.Errorf("FindByUserId userId: %s error: %v", in.UserId, err)
		return nil, err
	}
	return &pb.FindByIdResponse{
		UserId:   int64(user.Id),
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}
