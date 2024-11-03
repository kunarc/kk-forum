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
goctl model mysql datasource --dir ./internal/model --table user --cache true --url "root:8888@tcp(127.0.0.1:3306)/kk_user" --style go_zero
```
### 创建点赞主题并查看
```
cd /opt/kafka/bin
./kafka-topics.sh --create --topic kk-forum-like --bootstrap-server localhost:9092
./kafka-topics.sh --describe --topic kk-forum-like --bootstrap-server localhost:9092
```

### mysql创建canna用户
```
CREATE USER 'canal'@'%' IDENTIFIED WITH 'mysql_native_password' BY 'canal';
GRANT ALL PRIVILEGES ON *.* TO 'canal'@'%';
FLUSH PRIVILEGES;
```

### 检查kafka是否有canal传递的消息

```
在canal的配置文件中配置监听了kk_forum.like_count表， 我们通过改动该表的数据，观察kafka是否有消息消费
cd /opt/kafka/bin
./kafka-console-consumer.sh --topic kk-forum-like-count --bootstrap-server localhost:9092
```