package logic

import (
	"context"
	"pkg/jwt"
	"strings"

	"api/internal/code"
	"api/internal/grpc_client/user"
	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	if len(req.VerificationCode) == 0 {
		return nil, code.VerificationCodeEmpty
	}
	req.Mobile = strings.TrimSpace(req.Mobile)
	verCode, err := GetActiveVerCode(l.svcCtx.BizRedis, req.Mobile)
	if err != nil {
		if err == redis.Nil {
			return nil, code.VerificationCodeExpire
		}
		l.Logger.Errorf("Get active code error: err is %s", err.Error())
		return nil, err
	}
	if verCode != req.VerificationCode {
		return nil, code.VerificationCodeError
	}
	user, err := l.svcCtx.UserRpc.FindByMobile(l.ctx, &user.FindByMobileRequest{
		Mobile: req.Mobile,
	})
	if err != nil {
		l.Logger.Errorf("rpc find by mobile error: err is %s", err.Error())
		return nil, err
	}
	if user == nil {
		return nil, code.MobileIsError
	}
	token, err := jwt.BuildAccessToken(jwt.TokenOption{
		AccessSecretKey: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire:    l.svcCtx.Config.Auth.AccessExpire,
		Field: map[string]any{
			"userId": user.UserId,
		},
	})
	if err != nil {
		l.Logger.Errorf("gen user token error: err is %s", err.Error())
		return nil, err
	}
	return &types.LoginResponse{
		UserId: user.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AcessExpire,
		},
	}, nil
}
