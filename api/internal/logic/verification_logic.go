package logic

import (
	"context"
	"fmt"
	"pkg/util"
	"strconv"
	"strings"

	"api/internal/code"
	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	prefixVerCodeCount = "biz#verification#count#%s"
	prefixVerCodeKey   = "biz#verification#action#%s"
	sendVerCodeLimit   = 5
	verCodeExpries     = 60 * 2
	sendVerCodeExpries = 60 * 5
)

type VerificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerificationLogic {
	return &VerificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerificationLogic) Verification(req *types.VerificationRequest) (resp *types.VerificationResponse, err error) {
	// todo: add your logic here and delete this line
	mobile := strings.TrimSpace(req.Mobile)
	if len(mobile) != 11 {
		return nil, code.MobileIsError
	}
	count, err := l.getVerCodeCount(mobile)
	if err != nil {
		l.Logger.Errorf("Get Verification code error: mobile is %s, err is %s", mobile, err.Error())
	}
	if count >= sendVerCodeLimit {
		return nil, code.SendVerificationCodeExceed
	}
	verCode := util.RandomNumeric(6)
	// _, err = l.svcCtx.UserRpc.SendSms(l.ctx, &user.SendSmsRequest{Mobile: mobile})
	// if err != nil {
	// 	l.Logger.Errorf("send code to user error: mobile is %s, err is %s", mobile, err.Error())
	// 	return nil, code.SendVerCodeError
	// }
	err = l.saveVerCode(mobile, verCode)
	if err != nil {
		l.Logger.Errorf("save code to redis error: mobile is %s, err is %s", mobile, err.Error())
		return nil, code.SendVerCodeError
	}
	err = l.incrVerCode(mobile)
	if err != nil {
		l.Logger.Errorf("incr code count error: mobile is %s, err is %s", mobile, err.Error())
	}
	return &types.VerificationResponse{}, nil
}

// getVerCodeCount  获取验证码次数
//
//	@receiver l *VerificationLogic
//	@param mobile string
//	@return int
//	@return error
//	@author kunarc
//	@update 2024-10-25 07:48:49
func (l *VerificationLogic) getVerCodeCount(mobile string) (int, error) {
	key := fmt.Sprintf(prefixVerCodeCount, mobile)
	count, err := l.svcCtx.BizRedis.Get(key)
	if err != nil {
		return 0, nil
	}
	return strconv.Atoi(count)
}

func (l *VerificationLogic) saveVerCode(mobile, code string) error {
	key := fmt.Sprintf(prefixVerCodeKey, mobile)
	err := l.svcCtx.BizRedis.Setex(key, code, verCodeExpries)
	return err
}

func (l *VerificationLogic) incrVerCode(mobile string) error {
	key := fmt.Sprintf(prefixVerCodeCount, mobile)
	_, err := l.svcCtx.BizRedis.Incr(key)
	if err != nil {
		return err
	}
	return l.svcCtx.BizRedis.Expire(key, sendVerCodeExpries)
}

func GetActiveVerCode(rds *redis.Redis, mobile string) (string, error) {
	key := fmt.Sprintf(prefixVerCodeKey, mobile)
	return rds.Get(key)
}
