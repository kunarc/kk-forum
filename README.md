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
### 查看grpc服务列表/方法
```
grpcurl -plaintext 127.0.0.1:8080 list
grpcurl -plaintext 127.0.0.1:8080 list service.Like
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
### 创建ES索引 
```
PUT /article-index
{
    "mappings": {
      "properties": {
        "article_id": {
            "type": "long"
        },
        "title": {
            "type": "text"
        },
        "content": {
            "type": "text"
        },
        "description": {
            "type": "text"
        },
        "author_id": {
            "type": "long"
        },
        "author_name": {
            "type": "keyword"
        },
        "status": {
            "type": "integer"
        },
        "comment_num": {
            "type": "long"
        },
        "like_num": {
            "type": "long"
        },
        "collect_num": {
            "type": "long"
        },
        "view_num": {
            "type": "long"
        },
        "share_num": {
            "type": "long"
        },
        "tag_ids": {
            "type": "long"
        },
        "publish_time": {
            "type": "keyword"
        },
        "create_time": {
            "type": "keyword"
        },
        "update_time": {
            "type": "keyword"
        }
    }
  }
}
### 测试性能
```
ab -c 10 -n 11 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI3MzE4NDI4MDAsImlhdCI6MTczMTg0MjgwMCwidXNlcklkIjoxfQ.CVJyuEuLYXaDVK2oKcMj-dkYbzWaBkhAFt0DnrQ_tGY" "http://127.0.0.1:8888/v1/article/detail?article_id=1"
ab -c 10 -n 100000 "http://127.0.0.1:8888/v1/article/list?author_id=1&cursor=0&page_size=10&sort_type=0&article_id=0"
```