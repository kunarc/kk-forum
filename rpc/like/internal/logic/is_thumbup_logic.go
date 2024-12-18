package logic

import (
	"context"

	"like/internal/svc"
	"like/pb"

	// "like/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsThumbupLogic {
	return &IsThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsThumbupLogic) IsThumbup(in *pb.IsThumbupRequest) (*pb.IsThumbupResponse, error) {
	return &pb.IsThumbupResponse{}, nil
}
