# 星竟新棋牌项目web端

ETCD:配置中心

RabbitMQ:队列

Consul:服务注册发现

Redis:缓存

Mysql:数据库



RPC:微服务之间通信

HTTPS

WEBSOCKET
 
微服务(根据模块划分,根据业务划分)

协议：

![协议](https://github.com/HJM101/xj_web_server/blob/master/img/cmd.png)



//module 例子

模块文件夹

account: 用户模块,注册,登录,用户基本信息,token管理

apigw: 网关模块 需要Nginx负载,需要到通信模块注册（httpserver.routes.go）

//部署

ca module/account && ./start.sh

通信模块

cmd:websocket

部署命令：make && ./develop.sh

