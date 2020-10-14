module xj_web_server

go 1.13

require (
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.123
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/elazarl/goproxy v0.0.0-20191011121108-aa519ddbe484 // indirect
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.4.0
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/go-redis/redis/v7 v7.0.0-beta.4
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/xorm v0.7.9
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/websocket v1.4.1
	github.com/mattn/go-sqlite3 v1.11.0 // indirect
	github.com/micro/go-micro v1.14.0
	github.com/micro/go-plugins v1.4.0
	github.com/parnurzeal/gorequest v0.2.16
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271 // indirect
	google.golang.org/grpc v1.24.0
	gopkg.in/yaml.v2 v2.2.4
	moul.io/http2curl v1.0.0 // indirect
	xorm.io/core v0.7.2-0.20190928055935-90aeac8d08eb
)

replace (
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.1
	github.com/micro/go-plugins v1.3.0 => github.com/micro/go-plugins v1.3.0
)
