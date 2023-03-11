package controller

import (
	"YiHome/model"
	getCaptcha "YiHome/proto/getCaptcha"
	houseMicro "YiHome/proto/house"
	user "YiHome/proto/user"
	orderMicro "YiHome/proto/userOrder"
	"YiHome/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/guanguans/id-validator"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/satori/go.uuid"
	"image/png"
	"log"
	"net/http"
	"path/filepath"
)

//获取Session信息
func GetSession(c *gin.Context) {
	//resp := make(map[string]interface{})
	resp := Resp{}
	session := sessions.Default(c)
	userName := session.Get("userName")
	if userName == nil {
		resp.Errno = utils.RECODE_SESSIONERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_SESSIONERR)
		c.JSON(http.StatusOK, resp)
	} else {
		fmt.Println("获取到的session为：", userName)

		resp.Errno = utils.RECODE_OK
		resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
		//在Data结构体中传入map，为name字段
		tempMap := make(map[string]interface{})
		tempMap["name"] = userName.(string)
		resp.Data = tempMap

		c.JSON(http.StatusOK, resp)
	}
}

//获取验证码图片
func GetImageCode(c *gin.Context) {
	var img captcha.Image
	//获取传入的图片参数uuid
	uuid := c.Param("uuid")

	//初始化consul配置
	//consulReg := consul.NewRegistry( //registry.Addrs("139.9.208.92:8500"))
	//	func(options *registry.Options) {
	//		options.Addrs = []string{"47.115.229.57:8500"}
	//	})

	//创建一个microService,使用consul默认配置
	consulService := utils.NewMicroClient()

	//创建micro客户端，第一个参数为服务名称，第二个参数为上面创建的consulService的Client
	microClient := getCaptcha.NewGetCaptchaService("getCaptcha", consulService.Client())

	//var opts client.CallOption = func(o *client.CallOptions) {
	//	o.RequestTimeout = time.Second * 30
	//	o.DialTimeout = time.Second * 30
	//}
	//micro客户端调用call方法，使用resp接收返回的数据，
	//第一个参数为context.TODO( ) , 指定一个不为空的Context
	//第二个参数为request消息体
	resp, err := microClient.Call(context.TODO(), &getCaptcha.Request{Uuid: uuid})
	if err != nil {
		fmt.Println("microClient调用Call方法失败，错误为：", err)
		return
	}

	json.Unmarshal(resp.Img, &img)
	png.Encode(c.Writer, img)
	c.JSON(http.StatusOK, nil)
}

//调用微服务user.SendSms函数发送注册短信验证码
func SendMessege(c *gin.Context) {
	phone := c.Param("phone")
	text := c.Query("text")
	id := c.Query("id")

	////初始化consul配置
	//consulReg := consul.NewRegistry(func(options *registry.Options) {
	//	options.Addrs = []string{"47.115.229.57:8500"}
	//})
	////创建一个microService,使用consul默认配置
	//consulService := micro.NewService(
	//	micro.Registry(consulReg),
	//)
	consulService := utils.NewMicroClient()
	//创建micro客户端，第一个参数为服务名称，第二个参数为上面创建的consulService的Client
	microClient := user.NewUserService("user", consulService.Client())
	resp, err := microClient.SendSms(context.TODO(), &user.Request{Phone: phone, Uuid: id, ImgCode: text})
	if err != nil {
		fmt.Println("microClient调用SendSms方法失败，错误为：", err)
		c.JSON(http.StatusInternalServerError, "microClient调用SendSms方法失败")
		return
	}
	c.JSON(http.StatusOK, resp)
}

//注册用户信息
func PostRet(c *gin.Context) {
	var users User
	//创建User变量解析参数，如果解析错误则直接返回
	if err := c.ShouldBind(&users); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, fmt.Sprint("参数解析错误，", err))
		return
	}
	//consulReg := consul.NewRegistry(func(options *registry.Options) {
	//	options.Addrs = []string{"47.115.229.57:8500"}
	//})
	////创建一个microService,使用consul默认配置
	//consulService := micro.NewService(
	//	micro.Registry(consulReg),
	//)
	consulService := utils.NewMicroClient()
	//创建micro客户端，第一个参数为服务名称，第二个参数为上面创建的consulService的Client
	microClient := user.NewUserService("user", consulService.Client())
	fmt.Println("调用register方法")
	resp, err := microClient.Register(context.TODO(), &user.RegReq{Mobile: users.Mobile, SmsCode: users.Sms_Code, Password: users.PassWord})
	if err != nil {
		log.Println("microClient调用user.Register方法失败，错误为：", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

//获取地域信息
func GetArea(c *gin.Context) {
	var areas []model.Area

	coon := model.GetRedisPool()

	areaData, err := redis.Bytes(coon.Do("get", "areaData"))
	if err != nil {
		//如果Redis中获取不到，则从Mysql中获取序列化为json格式并写入Redis中
		model.DB().Debug().Find(&areas)

		areaBuf, _ := json.Marshal(areas) //序列化为json字节存入Redis
		_, err := coon.Do("Set", "areaData", areaBuf)
		if err != nil {
			log.Println("保存areas数据进入redis中出现错误，err为", err)
		}
	} else {
		//Redis中获取到数据则反序列化jinareas发送
		json.Unmarshal(areaData, &areas)
	}

	resp := Resp{
		Errno:  utils.RECODE_OK,
		Errmsg: utils.RecodeText(utils.RECODE_OK),
		Data:   areas,
	}

	c.JSON(http.StatusOK, resp)
}
func PostLogin(c *gin.Context) {
	var resp Resp
	var loginData struct {
		Mobile   string `json:"mobile" binding:"required"`
		PassWord string `json:"password" binding:"required"`
	}
	if err := c.ShouldBind(&loginData); err != nil {
		log.Println("loginData参数解析错误，err为：", err)
		c.JSON(http.StatusInternalServerError, fmt.Sprint("参数解析错误，", err))
	}
	fmt.Println(loginData)
	//验证手机号和密码是否正确，并返回用户名
	userName, err := model.Login(loginData.Mobile, loginData.PassWord)
	if err != nil {
		resp.Errno = utils.RECODE_LOGINERR
		resp.Errmsg = "用户登录失败,账号或密码错误"
		fmt.Println("账号密码错误")
		c.JSON(http.StatusOK, resp)
		return
	}
	fmt.Println("账号密码正确")
	//登录成功则将登录状态保存至session中
	session := sessions.Default(c)    //初始化session
	session.Set("userName", userName) //保存userName字段设置进session中
	session.Save()

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	c.JSON(http.StatusOK, resp)

}
func DeleteSessino(c *gin.Context) {
	resp := Resp{}
	session := sessions.Default(c) //初始化session
	session.Delete("userName")     //删除浏览器中的session
	err := session.Save()          //同时删除redis中的保存
	if err != nil {
		resp.Errno = "500"
		resp.Errmsg = "删除session失败"
	} else {
		resp.Errno = utils.RECODE_OK
		resp.Errmsg = utils.RECODE_OK
	}
	c.JSON(http.StatusOK, resp)
}
func GetUserInfo(c *gin.Context) {
	resp := Resp{}
	session := sessions.Default(c)
	userName := session.Get("userName")
	if userName == nil {
		resp.Errno = utils.RECODE_SESSIONERR
		resp.Errmsg = "获取用户名为空，非法登录界面，session为空"
		c.JSON(http.StatusOK, resp)
		return
	}
	userInfo, err := model.GetUserInfo(userName.(string))
	if err != nil {
		resp.Errno = utils.RECODE_SESSIONERR
		resp.Errmsg = "获取用户信息失败。"

		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	resp.Data = userInfo
	c.JSON(http.StatusOK, resp)
}
func PutUserInfo(c *gin.Context) {

	//获取用户名
	session := sessions.Default(c)
	userName := session.Get("userName")
	response := Resp{}
	//解析绑定传入的名字
	var nameData struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBind(&nameData); err != nil {
		log.Println("参数解析错误，err为：", err)
		c.JSON(http.StatusInternalServerError, fmt.Sprint("参数解析错误，err为：", err))
	}
	//修改用户名字
	err := model.UpdateUserName(nameData.Name, userName.(string))
	if err != nil {
		response.Errno = utils.RECODE_DBERR
		response.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		c.JSON(http.StatusOK, response)
		return
	}
	session.Set("userName", nameData.Name)
	err = session.Save()
	if err != nil {
		response.Errno = "500"
		response.Errmsg = "更新session错误"
		c.JSON(http.StatusOK, response)
		return
	}
	response.Errno = utils.RECODE_OK
	response.Errmsg = utils.RecodeText(utils.RECODE_OK)
	response.Data = nameData.Name
	c.JSON(http.StatusOK, response)
}

func PostAvatar(c *gin.Context) {
	//接收上传的头像
	file, err := c.FormFile("avatar")
	if err != nil {
		log.Println("上传图片出错，错误为：", err)
		c.JSON(http.StatusInternalServerError, "上传图片出错")
	}
	filename := uuid.NewV4()
	fmt.Println("生成的uuid为", filename)
	path := fmt.Sprintf("YiHome/avatar/%v", filename)
	savepath := path + filepath.Ext(file.Filename)
	token := utils.NewUpToken(path, savepath)

	if err := utils.PutFile(token, path, file); err != nil {
		log.Fatal("文件上传失败", err)
		c.JSON(http.StatusInternalServerError, "使用七牛云上传错误")
	}
	avatarUrl := storage.MakePublicURL("http://cdn.whyy1.top", savepath)

	userName := sessions.Default(c).Get("userName").(string)
	//err = model.UpdateAvatar(userName, avatarUrl)
	//if err != nil {
	//	log.Println("修改数据库头像地址错误，err为：", err)
	//	c.JSON(http.StatusInternalServerError, "头像修改失败")
	//	return
	//}

	//创建micro客户端，第一个参数为服务名称，第二个参数为上面创建的consulService的Client
	microClient := user.NewUserService("user", utils.NewMicroClient().Client())
	fmt.Println("调用register方法")
	response, err := microClient.UploadAvatar(context.TODO(), &user.UploadReq{
		UserName:  userName,
		AvatarUrl: avatarUrl,
	})
	if err != nil {
		log.Println("客户端调用user微服务UploadAvatar函数失败，err为：", err)
	}
	//存储在本地
	//c.SaveUploadedFile(file, "test/"+file.Filename)
	resp := Resp{
		Errno:  response.Errno,
		Errmsg: response.Errno,
		Data: map[string]string{
			"avatar_url": avatarUrl,
		},
	}
	c.JSON(http.StatusOK, resp)
}
func PostAuth(c *gin.Context) {
	var auth struct {
		IdCard   string `json:"id_card" binding:"required"`
		RealName string `json:"real_name" binding:"required,gte=2,lte=10,ne=1"`
	}
	if err := c.ShouldBind(&auth); err != nil {
		log.Println("实名验证请求参数解析错误，err为：", err)
		c.JSON(http.StatusOK, fmt.Sprint("实名验证请求参数解析错误，err为：", err))
		return
	}

	if result := idvalidator.IsValid(auth.IdCard, true); result != true {
		resp := Response{
			Errno:  utils.RECODE_DATAERR,
			Errmsg: utils.RecodeText("身份证信息有误，请仔细核对"),
		}
		c.JSON(http.StatusOK, resp)
		return
	}
	authInfo, _ := idvalidator.GetInfo(auth.IdCard, true)
	userName := sessions.Default(c).Get("userName").(string)

	//创建micro客户端，第一个参数为服务名称，第二个参数为上面创建的consulService的Client
	microClient := user.NewUserService("user", utils.NewMicroClient().Client())
	resp, err := microClient.AuthUpdate(context.TODO(), &user.AuthReq{
		IdCard:    auth.IdCard,
		RealName:  auth.RealName,
		UserName:  userName,
		IdAddress: authInfo.Address,
	})
	if err != nil {
		log.Println("客户端调用user微服务AuthUpdate函数失败，err为：", err)
	}
	fmt.Println(authInfo.Address, authInfo)

	c.JSON(http.StatusOK, resp)
}
func GetAuth(c *gin.Context) {
	resp := Resp{}
	userName := sessions.Default(c).Get("userName").(string)
	AuthInfo, _ := model.GetUserInfo(userName)
	fmt.Println(AuthInfo)
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	resp.Data = AuthInfo
	c.JSON(http.StatusOK, resp)
}

func GetOrders(ctx *gin.Context) {
	//获取get请求传参
	role := ctx.Query("role")
	//校验数据
	if role == "" {
		fmt.Println("获取数据失败")
		return
	}

	//处理数据  服务端处理业务
	microClient := orderMicro.NewUserorderService("userorder", utils.NewMicroClient().Client())
	//调用远程服务
	resp, _ := microClient.GetOrderInfo(context.TODO(), &orderMicro.GetReq{
		Role:     role,
		UserName: sessions.Default(ctx).Get("userName").(string),
	})

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}
func PostOrders(ctx *gin.Context) {
	var order struct {
		EndDate   string `json:"end_date"`
		HouseId   string `json:"house_id"`
		StartDate string `json:"start_date"`
	}
	if err := ctx.ShouldBind(&order); err != nil {
		fmt.Println(err)
		resp := Response{
			Errno:  utils.RECODE_DATAERR,
			Errmsg: utils.RecodeText(utils.RECODE_DATAERR),
		}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	//获取用户名
	userName := sessions.Default(ctx).Get("userName")

	//处理数据  服务端处理业务
	microClient := orderMicro.NewUserorderService("userorder", utils.NewMicroClient().Client())
	//调用服务
	resp, _ := microClient.CreateOrder(context.TODO(), &orderMicro.Request{
		StartDate: order.StartDate,
		EndDate:   order.EndDate,
		HouseId:   order.HouseId,
		UserName:  userName.(string),
	})

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}

//更新订单状态
func PutOrders(ctx *gin.Context) {
	//获取数据
	id := ctx.Param("id")
	var statusStu struct {
		Action string `json:"action"`
		Reason string `json:"reason"`
	}
	if err := ctx.ShouldBind(&statusStu); err != nil {
		fmt.Println(err)
		resp := Response{
			Errno:  utils.RECODE_DATAERR,
			Errmsg: utils.RecodeText(utils.RECODE_DATAERR),
		}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	//处理数据   更新订单状态
	microClient := orderMicro.NewUserorderService("userorder", utils.NewMicroClient().Client())
	//调用元和产能服务
	resp, _ := microClient.UpdateStatus(context.TODO(), &orderMicro.UpdateReq{
		Action: statusStu.Action,
		Reason: statusStu.Reason,
		Id:     id,
	})

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}
func PutOrdersComment(ctx *gin.Context) {

	//获取数据

	var commentStu struct {
		OrderId string `json:"order_id"`
		Comment string `json:"comment"`
	}
	fmt.Println(commentStu)

	if err := ctx.ShouldBind(&commentStu); err != nil {
		fmt.Println(err)
		resp := Response{
			Errno:  utils.RECODE_DATAERR,
			Errmsg: utils.RecodeText(utils.RECODE_DATAERR),
		}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	microClient := orderMicro.NewUserorderService("userorder", utils.NewMicroClient().Client())
	resp, _ := microClient.UpdateComment(context.TODO(), &orderMicro.CommentReq{
		OrderId: commentStu.OrderId,
		Comment: commentStu.Comment,
	})

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}

func GetIndex(ctx *gin.Context) {
	//处理数据
	microClient := houseMicro.NewHouseService("house", utils.NewMicroClient().Client())
	//调用服务
	resp, _ := microClient.GetIndexHouse(context.TODO(), &houseMicro.IndexReq{})

	ctx.JSON(http.StatusOK, resp)
}
