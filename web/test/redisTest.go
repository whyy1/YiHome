package main

import (
	"YiHome/conf"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func main() {
	config, err := conf.LoadConfig("../conf/")
	if err != nil {
		fmt.Println("Redis配置读取失败: ", err)
	}
	pool = &redis.Pool{ //实例化一个连接池
		MaxIdle: 16, //最初的连接数量
		// MaxActive:1000000,    //最大连接数量
		MaxActive: 0, //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		//IdleTimeout: 300, //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			return redis.Dial("tcp", config.REDIS_SOURCE)
		},
	}
	c := pool.Get()
	defer c.Close()
	fmt.Println("redis conn success")

	_, err = c.Do("Set", "qqq", "1006307055")
	if err != nil {
		fmt.Println("redis:设置token错误，id")
		return
	}
	s, err := redis.String(c.Do("Get", "s"))
	if err != nil {
		fmt.Println("redis:获取test错误，id", err)
		return
	}
	fmt.Println(s)
	defer c.Close()
}
