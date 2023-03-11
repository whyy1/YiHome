package utils

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"user/conf"
)

func NewAliyunClient() *dysmsapi.Client {
	config, err := conf.LoadConfig("../web/conf/")
	if err != nil {
		fmt.Println("读取阿里云短信配置信息失败: ", err)
	}
	credential := credentials.NewAccessKeyCredential(config.ACCESS_KEYID, config.ACCESS_KEYSECRET)

	client, err := dysmsapi.NewClientWithOptions("cn-hangzhou", sdk.NewConfig(), credential)
	if err != nil {
		panic(err)
	}
	return client
}
