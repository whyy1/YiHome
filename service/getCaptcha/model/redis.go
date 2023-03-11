package model

import (
	"getCaptcha/conf"
	"github.com/garyburd/redigo/redis"
	"log"
)

var pool *redis.Pool

func SaveImageCode(uuid, code string) error {
	log.Println("开始调用SaveImageCode函数")
	c := pool.Get()

	//开始将图片验证码设置到redis中
	if _, err := c.Do("setex", uuid, 5*60, code); err != nil {
		log.Println("redis:设置图片验证码错误，uuid为：", uuid, "  err为：", err)
		if err := c.Close(); err != nil {
			log.Println("图片验证码保存失败，SaveImageCode函数关闭redis连接时错误，err为：", err)
		}
		return err
	}

	log.Println("图片验证码保存成功")
	if err := c.Close(); err != nil {
		log.Println("图片验证码保存成功后，SaveImageCode函数关闭redis连接时错误，err为：", err)
	}
	log.Println("SaveImageCode函数调用结束")
	return nil
}
func init() {
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
