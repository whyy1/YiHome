package handler

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"log"
	"math/rand"
	"time"
	"user/model"
	"user/utils"

	pb "user/proto"
)

type User struct{}

func (e *User) SendSms(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	//判断验证码是否正确
	result := model.CheckImageCode(req.Uuid, req.ImgCode)
	//fmt.Println("result=", result)
	if result {
		//验证码输入正确则发送短信
		client := utils.NewAliyunClient()
		//创建发送信息Request
		request := dysmsapi.CreateSendSmsRequest()
		request.Scheme = "https"
		request.SignName = "韦海艺的个人网站"
		request.TemplateCode = "SMS_265005450"
		//生成6位数随机验证码
		rand.Seed(time.Now().UnixNano())
		smsCode := fmt.Sprintf("%06d", rand.Int31n(1000000))
		request.TemplateParam = fmt.Sprintf("{\"code\":\"%s\"}", smsCode)
		request.PhoneNumbers = req.Phone

		_, err := client.SendSms(request)
		if err != nil {
			fmt.Print(err.Error())
			rsp.Errno = utils.RECODE_SMSERR
			rsp.Errmsg = utils.RecodeText(utils.RECODE_SMSERR)
		}
		if err := model.SaveSmsCode(req.Phone, smsCode); err != nil {
			fmt.Println("Redis保存短信验证码失败，手机号为：", req.Phone, "  短信验证码为：", smsCode, "err为：", err)
		}
		rsp.Errno = utils.RECODE_OK
		rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	} else {
		//验证码错误则直接返回错误
		rsp.Errno = "500"
		rsp.Errmsg = utils.RecodeText("图形验证码错误。")
	}

	return nil
}
func (e *User) Register(ctx context.Context, req *pb.RegReq, rsp *pb.Response) error {
	//检验短信验证码，如果和redis中的不一样或不存在，直接返回错误
	if result := model.CheckSmsCode(req.Mobile, req.SmsCode); result != true {
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return nil
	}
	if err := model.RegisterUser(req.Mobile, req.Password); err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
	} else {
		rsp.Errno = utils.RECODE_OK
		rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	}
	return nil
}
func (e *User) UploadAvatar(ctx context.Context, in *pb.UploadReq, rsp *pb.Response) error {

	if err := model.UpdateAvatar(in.UserName, in.AvatarUrl); err != nil {
		log.Println("更新头像地址时失败，err为：", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
	} else {
		rsp.Errno = utils.RECODE_OK
		rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	}
	return nil
}
func (e *User) AuthUpdate(ctx context.Context, in *pb.AuthReq, rsp *pb.Response) error {

	if err := model.AuthUpdate(in.UserName, in.IdCard, in.RealName, in.IdAddress); err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
	} else {
		rsp.Errno = utils.RECODE_OK
		rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	}
	return nil
}
