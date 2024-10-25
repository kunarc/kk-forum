registry := "registry.cn-hangzhou.aliyuncs.com"
repository := "kunarc/kk-forum"
tag := "1.0"
arch_list := "linux/amd64,linux/arm64"
username := ""


# list all available recipes
default:
    @just --list --justfile {{ justfile() }}






docker-login:
    docker login --username={{ username }} registry.cn-hangzhou.aliyuncs.com  

docker:
    docker buildx build --no-cache -f ./Dockerfile --push --platform={{ arch_list }} -t {{registry}}/{{ repository }}:{{ tag }} .


# format source code
format:
    gofumpt -w cmd/ internal/


# tidy
tidy:
    go mod tidy

grpc:
    protoc --proto_path=pkg/xcode/pb --go_out=pkg/xcode/pb --go-grpc_out=pkg/xcode/pb  pkg/xcode/pb/status.proto
grpc-user:
    protoc --proto_path=api/internal/protos --go_out=api/internal/grpc_client --go-grpc_out=api/internal/grpc_client  api/internal/protos/user.proto