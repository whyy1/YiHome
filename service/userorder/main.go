package main

import (
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4/registry"
	"log"
	"userorder/conf"
	"userorder/handler"
	"userorder/model"
	pb "userorder/proto"

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

var (
	service = "userorder"
	version = "latest"
)

func main() {
	model.InitMysql()
	model.InitRedis()
	config, err := conf.LoadConfig("./conf/")
	if err != nil {
		log.Println("微服务house读取配置文件失败，err为", err)
	}
	//初始化consul配置
	var consulReg = consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{config.CONSUL_SOURCE} //指定要注册服务到哪个服务器的consul_service上，给出IP
	}) // Create service
	srv := micro.NewService()
	srv.Init(
		micro.Address(config.CONSUL_ADDRESS), //增加consul服务指定的地址
		micro.Name(config.CONSUL_SERVICE),    //服务名称，可以修改，但调用时需要一致
		micro.Version(config.CONSUL_VERSION),
		micro.Registry(consulReg), //使用consul作为服务发现
	)

	err = pb.RegisterUserorderHandler(srv.Server(), new(handler.Userorder))
	if err != nil {
		logger.Fatal(err)
		return
	}

	// Run service
	if err := srv.Run(); err != nil {
		log.Println("微服务house关闭！！")
		logger.Fatal(err)
	} else {
		log.Println("微服务house关闭！！")
	}
}
