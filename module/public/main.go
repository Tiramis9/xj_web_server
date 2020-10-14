package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	"github.com/micro/go-plugins/registry/etcdv3"
	_ "github.com/micro/go-plugins/transport/rabbitmq"
	"log"
	"xj_web_server/module/public/handler"
	"xj_web_server/module/public/proto"
)

func main() {
	eTCDRegistry := etcdv3.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:2379"}
	})
	service := micro.NewService(
		micro.Name("xj_web_server.service.public"),
		micro.Registry(eTCDRegistry),
		micro.Version("latest"),
	)
	// 初始化service, 解析命令行参数等
	service.Init()
	err := proto.RegisterPublicHandler(service.Server(), new(handler.Public))
	if err != nil {
		log.Fatalf("微服务注册失败:%v", err)
	}
	if err := service.Run(); err != nil {
		log.Fatalf("微服务启动失败:%v", err)
	}
}
