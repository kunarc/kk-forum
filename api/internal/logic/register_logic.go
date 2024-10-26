package logic

import (
	"context"
	"pkg/jwt"
	"pkg/xerror"
	"strings"

	"api/internal/code"
	"api/internal/grpc_client/user"
	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	// todo: add your logic here and delete this line
	name := req.Name
	req.Name = strings.TrimSpace(req.Name)
	if req.Name != name {
		return nil, code.RegisterUserNameSpace
	}
	if len(req.Name) == 0 {
		return nil, code.RegisterUserNameEmpty
	}
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
		return nil, xerror.ServerErr
	}
	if verCode != req.VerificationCode {
		return nil, code.VerificationCodeError
	}
	res, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterRequest{
		Username: req.Name,
		Mobile:   req.Mobile,
		Password: req.Password,
	})
	if err != nil {
		l.Logger.Errorf("call user rpc error: err is %s", err.Error())
		return nil, err
	}
	token, err := jwt.BuildAccessToken(jwt.TokenOption{
		AccessSecretKey: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire:    l.svcCtx.Config.Auth.AccessExpire,
		Field: map[string]any{
			"userId": res.UserId,
		},
	})
	if err != nil {
		l.Logger.Errorf("gen user token error: err is %s", err.Error())
		return nil, xerror.ServerErr
	}
	return &types.RegisterResponse{
		UserId: res.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AcessExpire,
		},
	}, nil
}
