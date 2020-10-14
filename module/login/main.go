package main

import (
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/transport/rabbitmq"
	"log"
	"xj_web_server/module/login/handler"
	"xj_web_server/module/login/proto"
	"xj_web_server/util"
)

func main() {
	service := micro.NewService(
		micro.Name("xj_web_server.service.login"),
		micro.Version("latest"),
	)

	// 初始化service, 解析命令行参数等
	service.Init(
		micro.BeforeStart(func() error {
			util.Logger.Infof("BeforeStart:%s","BeforeStart")
			return nil
		}),
		micro.AfterStart(func() error {
			util.Logger.Infof("AfterStart:%s","AfterStart")
			return nil
		}),
		micro.BeforeStop(func() error {
			util.Logger.Infof("BeforeStop:%s","BeforeStop")
			return nil
		}),
	)
	err := proto.RegisterLoginGRPCServiceHandler(service.Server(), new(handler.LoginGRPCService))
	if err != nil {
		log.Fatalf("微服务注册失败:%v", err)
	}

	if err := service.Run(); err != nil {
		log.Fatalf("微服务启动失败:%v", err)
	}
}