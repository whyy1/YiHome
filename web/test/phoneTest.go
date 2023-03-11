// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"YiHome/web/conf"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"math/rand"
	"time"
)

func sendSemessge(phone string) error {
	config, err := conf.LoadConfig("../conf/")
	if err != nil {
		fmt.Println("读取阿里云短信配置信息失败: ", err)
		return err
	}
	credential := credentials.NewAccessKeyCredential(config.ACCESS_KEYID, config.ACCESS_KEYSECRET)

	client, err := dysmsapi.NewClientWithOptions("cn-hangzhou", sdk.NewConfig(), credential)
	if err != nil {
		panic(err)
	}

	request := dysmsapi.CreateSendSmsRequest()

	request.Scheme = "https"

	request.SignName = "韦海艺的个人网站"
	request.TemplateCode = "SMS_265005450"
	request.TemplateParam = "{\"code\":\"133323\"}"
	fmt.Println("TemplateParam", request.TemplateParam)
	request.PhoneNumbers = phone

	response, err := client.SendSms(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
	return nil
}
func main() {
	//sendSemessge("13768660644")
	rand.Seed(time.Now().UnixNano())
	smsCode := fmt.Sprintf("%06d", rand.Int31n(1000000))
	fmt.Println("smsCode为", smsCode)
}
