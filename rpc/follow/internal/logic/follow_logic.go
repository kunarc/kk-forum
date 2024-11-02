package logic

import (
	"context"

	"follow/internal/svc"
	"follow/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 关注
func (l *FollowLogic) Follow(in *pb.FollowRequest) (*pb.FollowResponse, error) {
	return &pb.FollowResponse{}, nil
}
