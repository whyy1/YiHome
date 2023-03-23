package main

import (
	"YiHome/conf"
	"YiHome/controller"
	"YiHome/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	setupLogs()
	router := setupRouter()
	model.SetupRedis()
	model.SetupMysql()
	if err := router.RunTLS(":8080", "./conf/cert.crt", "./conf/cert.key"); err != nil {
		log.Println("gin框架运行失败，err为：", err)
		return
	}
}
func setupRouter() *gin.Engine {
	router := gin.Default()

	router.Static("/home", "view")
	config, err := conf.LoadConfig("./conf/")
	if err != nil {
		log.Println("读取Redis配置文件失败: ", err)
	}
	//初始化容器
	store, err := redis.NewStore(10, "tcp", config.REDIS_SOURCE, config.REDIS_PASSWORD, []byte("whyy1"))
	if err != nil {
		log.Println("初始化session容器失败，错误为：", err)
	}

	//使用容器
	router.Use(sessions.Sessions("mysession", store))

	r1 := router.Group("/api/v1.0/")
	{
		r1.GET("session", controller.GetSession)
		r1.GET("imagecode/:uuid", controller.GetImageCode)
		r1.GET("smscode/:phone", controller.SendMessege)
		r1.POST("users", controller.PostRet)
		r1.GET("areas", controller.GetArea)
		r1.POST("sessions", controller.PostLogin)

		r1.Use(LoginFilter()) //使用中间件检查session是否存在，不存在直接返回

		//删除session
		r1.DELETE("session", controller.DeleteSessino)
		//获取用户信息
		r1.GET("user", controller.GetUserInfo)
		//上传用户头像
		r1.POST("user/avatar", controller.PostAvatar)
		//修改用户名称
		r1.PUT("user/name", controller.PutUserInfo)

		//获取实名认证信息
		r1.GET("user/auth", controller.GetAuth)
		//上传实名认证信息
		r1.POST("user/auth", controller.PostAuth)

		//获取网站首页图和连接
		r1.GET("house/index", controller.GetIndex)
		//搜索发布的房屋信息
		r1.GET("houses", controller.GetHouses)
		r1.GET("user/houses", controller.GetUserHouses)
		//获取房屋的详细信息
		r1.GET("houses/:id", controller.GetHousesInfo)
		//上传房源信息
		r1.POST("houses", controller.PostHouses)
		//上传房屋图片
		r1.POST("houses/:id/images", controller.PostHousesImage)
		//修改房屋的详细信息
		r1.PUT("houses/:id", controller.PutHousesInfo)
		//删除房屋的详细信息
		r1.DELETE("houses/:id", controller.DeleteHousesInfo)
		//获取用户发布的房屋信息

		//获取订单
		r1.GET("user/orders", controller.GetOrders)
		//客户预订房屋订单
		r1.POST("orders", controller.PostOrders)
		//房东修改订单状态
		r1.PUT("orders/:id/status", controller.PutOrders)
		//租客提交订单评论
		r1.PUT("orders/:id/comment", controller.PutOrdersComment)

	}

	return router
}
func setupLogs() {
	logFile, err := os.OpenFile("./logs/web.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	log.Println("日志打开成功！")
}
func LoginFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userName := session.Get("userName")
		if userName == nil {
			c.Abort()
		} else {
			c.Next()
		}
	}
}
