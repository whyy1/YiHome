package main

import (
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4/registry"
	"log"
	"user/conf"
	"user/handler"
	"user/model"
	pb "user/proto"

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

var (
	service = "user"
	version = "latest"
)

func main() {
	model.SetupMysql()
	model.SetupRedis()
	config, err := conf.LoadConfig("./conf/")
	if err != nil {
		log.Println("微服务user读取配置文件失败，err为", err)
	}
	//初始化consul配置
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{config.CONSUL_SOURCE} //指定要注册服务到哪个服务器的consul_service上，给出IP
	})
	// Create service
	srv := micro.NewService()
	srv.Init(
		micro.Address(config.CONSUL_ADDRESS), //增加consul服务指定的地址
		micro.Name(config.CONSUL_SERVICE),    //服务名称，可以修改，但调用时需要一致
		micro.Version(config.CONSUL_VERSION),
		micro.Registry(consulReg), //使用consul作为服务发现
	)

	// Register handler
	if err := pb.RegisterUserHandler(srv.Server(), new(handler.User)); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	} else {
		log.Println("微服务user启动成功！！")
	}
}
