package main

import (
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/encoder/json"
	"github.com/micro/go-micro/config/source"
	"github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	"github.com/micro/go-plugins/config/source/consul"
	_ "github.com/micro/go-plugins/transport/rabbitmq"
	"log"
	"xj_web_server/module/account/handler"
	"xj_web_server/module/account/proto"
	"time"
)

func main() {
	consulSource := consul.NewSource(
		// optionally specify etcd address; default to localhost:8500
		consul.WithAddress("localhost:8500"),
		// optionally specify prefix; defaults to /micro/config
		consul.WithPrefix("/my/prefix"),
		// optionally strip the provided prefix from the keys, defaults to false
		consul.StripPrefix(true),
		source.WithEncoder(json.NewEncoder()),
	)
	// Create new config
	conf := config.NewConfig()

	// Load file source
	err := conf.Load(consulSource)

	if err != nil {
		log.Fatalf("微服务配置加载失败:%v", err)
	}

	fmt.Println(string(conf.Get("1").Bytes()))

	//将配置写入结构
	type Host struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	}

	type Config struct {
		Hosts map[string]Host `json:"hosts"`
	}
	var configVar Config
	err = conf.Get("1").Scan(&configVar)
	// consul kv get my/prefix/host
	fmt.Println(configVar, err)

	//读取多个值
	//如果将配置写入结构
	var host Host
	err = conf.Get("1", "hosts", "database").Scan(&host)
	fmt.Println(host, err)

	//读取独立的值
	// 获取address值，缺省值使用 “localhost”
	address := conf.Get("1", "hosts", "database", "address").String("localhost")
	// 获取port值，缺省值使用 3000
	port := conf.Get("1", "hosts", "database", "port").Int(3000)
	fmt.Println(address, port)

	go func() {
		for {
			//观测目录的变化。当文件有改动时，新值便可生效。
			w, err := conf.Watch("1", "hosts", "database")
			if err != nil {
				// do something
				log.Fatalf("微服务配置watch失败:%v", err)
			}
			// wait for next value
			v, err := w.Next()
			if err != nil {
				// do something
				log.Fatalf("微服务配置watch失败:%v", err)
			}
			err = v.Scan(&host)
			if err != nil {
				// do something
				log.Fatalf("微服务配置watch失败:%v", err)
			}
			fmt.Println(host)
		}
	}()

	//----代理设置位置-----
	//cmd.Init()
	//if err := broker.Init(); err != nil {
	//	log.Fatalf("Broker Init error: %v", err)
	//}
	//if err := broker.Connect(); err != nil {
	//	log.Fatalf("Broker Connect error: %v", err)
	//}
	b := rabbitmq.NewBroker(
		broker.Addrs("amqp://guest:guest@127.0.0.1:5672"),
	)
	if err := b.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}
	if err := b.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}


	//go pub(b)
	go sub(b)
	//-----代理设置位置end -----
	service := micro.NewService(
		micro.Name("xj_web_server.service.user"),
		micro.Version("latest"),
		micro.Broker(b),
	)

	// 初始化service, 解析命令行参数等
	service.Init(
		micro.BeforeStart(func() error {
			//util.Logger.Infof("BeforeStart:%s","BeforeStart")
			return nil
		}),
		micro.AfterStart(func() error {
			//util.Logger.Infof("AfterStart:%s","AfterStart")
			return nil
		}),
		micro.BeforeStop(func() error {
			//util.Logger.Infof("BeforeStop:%s","BeforeStop")
			return nil
		}),
	)
	err = proto.RegisterUserHandler(service.Server(), new(handler.User))
	if err != nil {
		log.Fatalf("微服务注册失败:%v", err)
	}

	if err := service.Run(); err != nil {
		log.Fatalf("微服务启动失败:%v", err)
	}
}

var (
	topic = "go.micro.topic.foo"
)

func pub(b broker.Broker) {
	tick := time.NewTicker(time.Second)
	i := 0
	for _ = range tick.C {
		msg := &broker.Message{
			Header: map[string]string{
				"id": fmt.Sprintf("%d", i),
			},
			Body: []byte(fmt.Sprintf("%d: %s", i, time.Now().String())),
		}
		if err := b.Publish(topic, msg); err != nil {
			log.Printf("[pub] failed: %v", err)
		} else {
			fmt.Println("[pub] pubbed message:", string(msg.Body))
		}
		i++
	}
}

func sub(b broker.Broker) {
	_, err := b.Subscribe(topic, func(p broker.Event) error {

		fmt.Println("[sub] received message:", string(p.Message().Body), "header", p.Message().Header, "topic", p.Topic())
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
