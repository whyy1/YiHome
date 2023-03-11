package main

import (
	"github.com/afocus/captcha"
	"image/color"
	"image/png"
	"net/http"
)

func main() {
	cap := captcha.New()
	//设置字体,一定要指定文件路径
	cap.SetFont("C:\\Users\\why\\go\\pkg\\mod\\github.com\\afocus\\captcha@v0.0.0-20191010092841-4bd1f21c8868\\examples\\comic.ttf")
	// 设置验证码大小
	cap.SetSize(128, 64)
	// 设置干扰强度
	cap.SetDisturbance(captcha.MEDIUM)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	cap.SetFrontColor(color.White)
	// 设置背景色 可以多个 随机替换背景色 默认白色
	cap.SetBkgColor(color.Black)

	http.HandleFunc("/r", func(w http.ResponseWriter, r *http.Request) {
		//创建验证码6个字符，captcha.NUM字符模式数字类型
		//返回验证码图像对象
		img, str := cap.Create(6, captcha.NUM)
		png.Encode(w, img)
		println(str)
	})

	http.ListenAndServe(":8085", nil)
}
