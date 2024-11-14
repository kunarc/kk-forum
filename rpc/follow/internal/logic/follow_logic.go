package logic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"follow/code"
	"follow/internal/model"
	"follow/internal/svc"
	"follow/internal/types"
	"follow/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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
	if in.UserId == 0 {
		return nil, code.FollowUserIdEmpty
	}
	if in.FollowedUserId == 0 {
		return nil, code.FollowedUserIdEmpty
	}
	if in.UserId == in.FollowedUserId {
		return nil, code.CannotFollowSelf
	}
	follow, err := l.svcCtx.FollowModel.FindByUserIDAndFollowedUserID(l.ctx, in.UserId, in.FollowedUserId)
	if err != nil {
		l.Logger.Errorf("[Follow] FollowModel.FindByUserIDAndFollowedUserID err: %v req: %v", err, in)
		return nil, err
	}
	if follow != nil && follow.FollowStatus == types.FollowStatusFollow {
		return &pb.FollowResponse{}, nil
	}
	followTime := time.Now()
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if follow != nil {
			// 更新follow状态
			err = model.NewFollowModel(tx).UpdateFollowStatus(l.ctx, followTime, follow.ID, types.FollowStatusFollow)
			if err != nil {
				l.Logger.Errorf("[UpdateFollowStatus] error: err is %v, id is %v", err.Error(), follow.ID)
			}
		} else {
			follow := &model.Follow{
				UserID:         in.UserId,
				FollowedUserID: in.FollowedUserId,
				FollowStatus:   types.FollowStatusFollow,
				CreateTime:     followTime,
				UpdateTime:     followTime,
			}
			err = model.NewFollowModel(tx).InsertFollowRecord(l.ctx, follow)
			if err != nil {
				l.Logger.Errorf("[InsertFollowRecord] error: err is %v", err.Error())
			}
			err = model.NewFollowCountModel(tx).InsertFollowCount(l.ctx, in.UserId, in.FollowedUserId, followTime)
			if err != nil {
				l.Logger.Errorf("[InsertFollowCount] error: err is %v", err.Error())
			}
		}
		// 更新 follow_count 表
		err = model.NewFollowCountModel(tx).IncrFollowCount(l.ctx, in.UserId, followTime)
		if err != nil {
			l.Logger.Errorf("[IncrFollowCount] error: err is %v, userId is %v", err.Error(), in.UserId)
		}
		err = model.NewFollowCountModel(tx).IncrFansCount(l.ctx, in.FollowedUserId, followTime)
		if err != nil {
			l.Logger.Errorf("[IncrFansCount] error: err is %v, userId is %v", err.Error(), in.FollowedUserId)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	followTimeInt64, _ := strconv.ParseInt(followTime.Format("2006-01-02 15:04:05"), 10, 64)
	followExist, _ := l.svcCtx.BizRedis.ExistsCtx(l.ctx, cacheFollowKey(types.CacheFollowType, in.UserId))
	if followExist {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, cacheFollowKey(types.CacheFollowType, in.UserId), followTimeInt64, strconv.FormatInt(in.FollowedUserId, 10))
		if err != nil {
			l.Logger.Errorf("[ZaddCtx] cache follow key error: err is %v, key is %v", err.Error(), cacheFollowKey(types.CacheFollowType, in.UserId))
		}
	}
	fansExist, _ := l.svcCtx.BizRedis.ExistsCtx(l.ctx, cacheFollowKey(types.CacheFansType, in.UserId))
	if fansExist {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, cacheFollowKey(types.CacheFollowType, in.FollowedUserId), followTimeInt64, strconv.FormatInt(in.UserId, 10))
		if err != nil {
			l.Logger.Errorf("[ZaddCtx] cache fans key error: err is %v, key is %v", err.Error(), cacheFollowKey(types.CacheFansType, in.FollowedUserId))
		}
	}
	return &pb.FollowResponse{}, nil
}

func cacheFollowKey(cacheType string, uID int64) (key string) {
	if cacheType == types.CacheFollowType {
		key = fmt.Sprintf("%s%d", types.CacheFollowPrefix, uID)
	} else if cacheType == types.CacheFansType {
		key = fmt.Sprintf("%s%d", types.CacheFansPrefix, uID)
	}
	return
}
