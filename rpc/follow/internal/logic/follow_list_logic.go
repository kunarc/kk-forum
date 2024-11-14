package logic

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"follow/code"
	"follow/internal/svc"
	"follow/internal/types"
	"follow/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/threading"
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
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor == 0 {
		in.Cursor = time.Now().Unix()
	}

	var (
		isCache        = false
		followItems    []*pb.FollowItem
		lastId, cursor int64
		err            error
	)
	followerIds, _ := l.cacheFollowOrFansIDs(in.UserId, in.Cursor, in.PageSize, types.CacheFollowType)
	// if err != nil {
	// 	l.Logger.Errorf("[cacheFollowOrFansIDs] is error: err is %v, cacheType is %v", err.Error(), types.CacheFollowPrefix)
	// 	return nil, err
	// }
	if len(followerIds) > 0 {
		isCache = true
		followItems, _ = l.FollowItemByFollowedIds(l.ctx, in.UserId, followerIds)

	} else {
		followItems, err = l.svcCtx.FollowModel.GetFollowItemList(l.ctx, in.UserId, in.Cursor, in.PageSize)
		if err != nil {
			l.Logger.Errorf("[FollowModel.GetFollowItemList] is error: err is %v", err.Error())
			return nil, err
		}
	}
	if len(followItems) > 0 {
		lastId, cursor = followItems[len(followItems)-1].Id, followItems[len(followItems)-1].CreateTime
		for k, ft := range followItems {
			if ft.Id == in.Id && ft.CreateTime == in.Cursor {
				followItems = followItems[k:]
				break
			}
		}
	}
	if !isCache {
		threading.GoSafe(func() {
			if len(followItems) > 0 {
				err := l.addCacheFollowedUserId(in.UserId, followItems)
				if err != nil {
					l.Logger.Errorf("[addCacheFollowedUserId] is error: err is %v, uid is %v", err.Error(), in.UserId)
				}
			}
		})
	}
	return &pb.FollowListResponse{
		Items:  followItems,
		Cursor: cursor,
		Id:     lastId,
	}, nil
}

// 获取缓存信息
func (l *FollowListLogic) cacheFollowOrFansIDs(uID, cursor, ps int64, keyType string) ([]int64, error) {
	var key string
	if keyType == types.CacheFollowType {
		key = fmt.Sprintf("%s%d", types.CacheFollowPrefix, uID)
	} else if keyType == types.CacheFansType {
		key = fmt.Sprintf("%s%d", types.CacheFansPrefix, uID)
	}
	var followIDs []int64
	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(l.ctx, key, 0, cursor, 0, int(ps))
	if err != nil {
		l.Logger.Errorf("[ZrevrangebyscoreWithScoresAndLimitCtx] get cache error: err is %v, key is %v, cursor is %v", err.Error(), key, cursor)
		return nil, err
	}
	b, _ := l.svcCtx.BizRedis.ExistsCtx(l.ctx, key)
	if b {
		err := l.svcCtx.BizRedis.ExpireCtx(l.ctx, key, types.FollowKeyExpire)
		if err != nil {
			return nil, err
		}
	}
	for _, pair := range pairs {
		id, _ := strconv.ParseInt(pair.Key, 10, 64)
		followIDs = append(followIDs, id)
	}
	return followIDs, nil
}

func (l *FollowListLogic) FollowItemByFollowedIds(ctx context.Context, uid int64, followedIds []int64) ([]*pb.FollowItem, error) {
	followItems, err := mr.MapReduce(func(source chan<- int64) {
		for _, fid := range followedIds {
			source <- fid
		}
	}, func(fid int64, writer mr.Writer[*pb.FollowItem], cancel func(error)) {
		f, err := l.svcCtx.FollowModel.FindOne(l.ctx, uid, fid)
		if err != nil {
			l.Logger.Errorf("[FollowModel.FindOne] is error: err is %v", err.Error())
			cancel(err)
		}
		fc, err := l.svcCtx.FollowCountModel.FindOne(ctx, fid)
		if err != nil {
			l.Logger.Errorf("[FollowCountModel.FindOne] is error: err is %v", err.Error())
			cancel(err)
		}
		p := &pb.FollowItem{
			Id:             f.ID,
			FollowedUserId: fid,
			FansCount:      int64(fc.FansCount),
			CreateTime:     f.CreateTime.Unix(),
		}
		writer.Write(p)
	}, func(pipe <-chan *pb.FollowItem, writer mr.Writer[[]*pb.FollowItem], cancel func(error)) {
		var followItems []*pb.FollowItem
		for p := range pipe {
			followItems = append(followItems, p)
		}
		// 排序
		sort.Slice(followItems, func(i, j int) bool {
			return followItems[i].CreateTime > followItems[j].CreateTime
		})
		writer.Write(followItems)
	})
	if err != nil {
		return nil, err
	}
	return followItems, nil
}

func (l *FollowListLogic) addCacheFollowedUserId(uID int64, fItem []*pb.FollowItem) error {
	key := fmt.Sprintf("%s%d", types.CacheFollowPrefix, uID)
	var pairs []redis.Pair
	for _, item := range fItem {
		pairs = append(pairs, redis.Pair{
			Key:   strconv.FormatInt(item.FollowedUserId, 10),
			Score: item.CreateTime,
		})
	}
	if _, err := l.svcCtx.BizRedis.Zadds(key, pairs...); err != nil {
		return err
	}
	return l.svcCtx.BizRedis.Expire(key, types.FollowKeyExpire)
}
