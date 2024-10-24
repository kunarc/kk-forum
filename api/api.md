### 1. N/A

1. route definition

- Url: /v1/login
- Method: POST
- Request: `LoginRequest`
- Response: `LoginResponse`

2. request definition



```golang
type LoginRequest struct {
	Mobile string `json:"mobile"`
	VerificationCode string `json:"verification_code"`
}
```


3. response definition



```golang
type LoginResponse struct {
	UserId int64 `json:"userId"`
	Token Token `json:"token"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	AccessExpire int64 `json:"access_expire"`
}
```

### 2. N/A

1. route definition

- Url: /v1/register
- Method: POST
- Request: `RegisterRequest`
- Response: `RegisterResponse`

2. request definition



```golang
type RegisterRequest struct {
	Name string `json:"name"`
	Mobile string `json:"mobile"`
	Password string `json:"password"`
	VerificationCode string `json:"verification_code"`
}
```


3. response definition



```golang
type RegisterResponse struct {
	UserId int64 `json:"user_id"`
	Token Token `json:"token"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	AccessExpire int64 `json:"access_expire"`
}
```

### 3. N/A

1. route definition

- Url: /v1/verification
- Method: POST
- Request: `VerificationRequest`
- Response: `VerificationResponse`

2. request definition



```golang
type VerificationRequest struct {
	Mobile string `json:"mobile"`
}
```


3. response definition



```golang
type VerificationResponse struct {
}
```

### 4. N/A

1. route definition

- Url: /v1/user/info
- Method: GET
- Request: `-`
- Response: `UserInfoResponse`

2. request definition



3. response definition



```golang
type UserInfoResponse struct {
	UserId int64 `json:"user_id"`
	Username string `json:"username"`
	Avatar string `json:"avatar"`
}
```

