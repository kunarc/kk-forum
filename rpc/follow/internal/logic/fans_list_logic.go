package logic

import (
	"context"

	"follow/internal/svc"
	"follow/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FansListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFansListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FansListLogic {
	return &FansListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 粉丝列表
func (l *FansListLogic) FansList(in *pb.FansListRequest) (*pb.FansListResponse, error) {
	return &pb.FansListResponse{}, nil
}
