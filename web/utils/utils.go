package utils

import (
	"YiHome/conf"
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"log"
)

func NewMicroClient() micro.Service {
	config, err := conf.LoadConfig(".")
	if err != nil {
		log.Fatal("初始化Micro客户端出错，读取ConsulAdress失败，错误为：", err)
	}
	//初始化consul配置
	consulReg := consul.NewRegistry(
		func(options *registry.Options) {
			options.Addrs = []string{config.Consul_Address}
		})
	//创建一个microService,使用consul默认配置
	consulService := micro.NewService(
		//micro.Address(config.CONSUL_ADDRESS), //增加consul服务指定的地址
		//micro.Name("getCaptcha"),             //服务名称，可以修改，但调用时需要一致
		//micro.Version("latest"),
		micro.Registry(consulReg), //使用consul作为服务发现
	)
	return consulService
}
