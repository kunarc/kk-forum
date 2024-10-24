# kk-forum
微服务论坛后台

### 生成api服务

```
cd api
goctl api go --dir ./ --api api.api --style go_zero
```
### 生成rpc服务（用户服务为例）
```
cd rpc/user
goctl rpc protoc ./user.proto --go_out=. --go-grpc_out=. --zrpc_out ./ --style go_zero
```
### 生成模型
```
cd rpc/user
goctl model mysql datasource --dir ./internal/model --table user --cache true --url "root:200483@tcp(127.0.0.1:3306)/beyond_user" --style go_zero
```