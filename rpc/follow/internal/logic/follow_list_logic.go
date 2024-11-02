package logic

import (
	"context"

	"follow/internal/svc"
	"follow/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 关注列表
func (l *FollowListLogic) FollowList(in *pb.FollowListRequest) (*pb.FollowListResponse, error) {
	return &pb.FollowListResponse{}, nil
}
