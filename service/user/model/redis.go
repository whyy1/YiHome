package model

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"user/conf"
)

var pool *redis.Pool

func CheckImageCode(uuid, imgCode string) bool {
	c := pool.Get()
	defer c.Close()

	code, err := redis.String(c.Do("Get", uuid))
	if err != nil {
		fmt.Println("redis查询图片验证码错误，uuid为：", uuid, "  err为", err)
		return false
	}
	return code == imgCode
}
func SaveSmsCode(phone, code string) error {
	c := pool.Get()
	defer c.Close()

	code, err := redis.String(c.Do("setex", phone, 5*60, code))
	if err != nil {
		return err
	}
	return nil
}
func CheckSmsCode(phone, code string) bool {
	c := pool.Get()
	defer c.Close()
	smscode, err := redis.String(c.Do("Get", phone))
	if err != nil {
		fmt.Printf("redis【查询】短信验证码错误，phone为：%verr为:%v\n", phone, err)
		return false
	}
	return smscode == code
}

func SetupRedis() {
	config, err := conf.LoadConfig("./conf/")
	if err != nil {
		log.Println("初始化Redis失败，读取Redis配置文件失败: ", err)
		return
	}

	pool = &redis.Pool{ //实例化一个连接池
		MaxIdle: 16, //最初的连接数量
		// MaxActive:1000000,    //最大连接数量
		MaxActive: 0, //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		//IdleTimeout: 300, //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			return redis.Dial("tcp", config.REDIS_SOURCE, redis.DialPassword(config.REDIS_PASSWORD))
		},
	}
	log.Println("初始化Redis成功")
}
