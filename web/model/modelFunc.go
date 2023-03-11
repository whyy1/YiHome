package model

import (
	"YiHome/conf"
	"crypto/md5"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
)

var pool *redis.Pool

func Login(mobile, password string) (string, error) {
	var user User

	pwd_hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	err := DB().Where("mobile = ? AND password_hash = ?", mobile, pwd_hash).
		First(&user).Error
	if err != nil {
		log.Printf("用户登录失败！手机号为：%v,密码为:%v,err:为%v\n", mobile, password, err)
	}
	return user.Name, err
}
func GetUserInfo(name string) (User, error) {
	var user User

	err := DB().Where("name = ?", name).First(&user).Error
	if err != nil {
		log.Printf("获取用户信息失败！用户名为：%v,   err:为%v\n", name, err)
		return User{}, err
	}
	user.Password_hash = ""
	return user, err
}
func UpdateUserName(name, oldName string) error {

	return DB().Model(User{}).Where("name=?", oldName).Update("name", name).Error
}
func UpdateAvatar(userName, avatar string) error {

	return DB().Model(User{}).Where("name=?", userName).Update("avatar_url", avatar).Error
}

func GetAuthInfo(name string) (User, error) {
	user := User{}
	err := DB().Where("name = ?", name).First(&user).Error
	if err != nil {
		log.Printf("获取用户信息失败！用户名为：%v,   err:为%v\n", name, err)
		return User{}, err
	}
	return user, nil
}
func GetOrders(name string) {

}
func GetRedisPool() redis.Conn {
	return pool.Get()
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
