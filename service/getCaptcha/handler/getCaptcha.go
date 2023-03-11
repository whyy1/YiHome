package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"getCaptcha/model"
	pb "getCaptcha/proto"
	"github.com/afocus/captcha"
	"image/color"
	"log"
)

type GetCaptcha struct{}

func (e *GetCaptcha) Call(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	log.Println("开始调用getCaptcha微服务")
	//生成随机图片验证码
	cap := captcha.New()
	//设置字体,一定要指定文件路径
	if err := cap.SetFont("conf/comic.ttf"); err != nil {
		fmt.Println("验证码字体文件读取失败，err为：", err)
		return err
	}
	// 设置验证码大小
	cap.SetSize(128, 64)
	// 设置干扰强度
	cap.SetDisturbance(captcha.MEDIUM)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	cap.SetFrontColor(color.White)
	// 设置背景色 可以多个 随机替换背景色 默认白色
	cap.SetBkgColor(color.Black)
	//创建图片验证码和字符串
	img, str := cap.Create(6, captcha.NUM)
	//转换成bytes接收，
	imgBuf, _ := json.Marshal(img)
	rsp.Img = imgBuf

	return model.SaveImageCode(req.Uuid, str)
}
