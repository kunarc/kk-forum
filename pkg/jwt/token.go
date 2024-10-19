package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (

	// TokenOption struct
	//	@update 2024-10-19 04:45:51
	TokenOption struct {
		AccessSecretKey string
		AccessExpire    int64
		Field           map[string]any
	}

	// Token struct
	//	@update 2024-10-19 04:58:11
	Token struct {
		AccessToken string `json:"access_token"`
		AcessExpire int64  `json:"acess_expire"`
	}
)

// BuildAccessToken
//
//	@param opt TokenOption
//	@return *Token
//	@return error
//	@author kunarc
//	@update 2024-10-19 05:16:42
func BuildAccessToken(opt TokenOption) (*Token, error) {
	now := time.Now().Add(-time.Minute).Unix()
	accessToken, err := genAccessToken(now, opt.AccessExpire, opt.AccessSecretKey, opt.Field)
	if err != nil {
		return nil, err
	}
	return &Token{
		AccessToken: accessToken,
		AcessExpire: now + opt.AccessExpire,
	}, nil
}

// genAccessToken
//
//	@param iat int64
//	@param expire int64
//	@param secretKey string
//	@param filed map[string]any
//	@return string
//	@return error
//	@author kunarc
//	@update 2024-10-19 05:12:04
func genAccessToken(iat int64, expire int64, secretKey string, filed map[string]any) (string, error) {
	chaims := make(jwt.MapClaims)
	{
		chaims["iat"] = iat
		chaims["exp"] = iat + expire
	}
	for k, v := range filed {
		chaims[k] = v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, chaims)
	return token.SignedString([]byte(secretKey))
}
