package main

import (
	"getCaptcha/conf"
	"getCaptcha/handler"
	pb "getCaptcha/proto"
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"log"
)

func main() {
	config, err := conf.LoadConfig("./conf/")
	if err != nil {
		log.Println("微服务user读取配置文件失败，err为", err)
	}
	//初始化consul配置
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{config.CONSUL_SOURCE} //指定要注册服务到哪个服务器的consul_service上，给出IP
	})

	//指定提供微服务的IP和端口号，
	srv := micro.NewService()
	srv.Init(
		micro.Address(config.CONSUL_ADDRESS), //指定提供微服务的地址
		micro.Name(config.CONSUL_SERVICE),    //微服务名称，可以修改，但调用时需要一致
		micro.Version(config.CONSUL_VERSION), //微服务版本
		micro.Registry(consulReg),            //使用consul作为服务发现，将微服务注册到consul上
	)

	// Register handler
	if err := pb.RegisterGetCaptchaHandler(srv.Server(), new(handler.GetCaptcha)); err != nil {
		logger.Fatal(err)
	}
	// 启动微服务
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	} else {
		log.Println("微服务getCaptcha启动成功！！")
	}
}
