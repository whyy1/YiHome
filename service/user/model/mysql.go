package model

import (
	"crypto/md5"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"user/conf"
)

var db *gorm.DB

func RegisterUser(Mobile, pwd string) error {
	//将密码进行md5加密
	password := fmt.Sprintf("%x", md5.Sum([]byte(pwd)))
	//讲手机号和加密后的密码存入mysql
	user := User{
		Name:          Mobile,
		Password_hash: password,
		Mobile:        Mobile,
	}
	fmt.Println("mysql保存用户")
	return db.Debug().Create(&user).Error
}
func UpdateAvatar(userName, avatar string) error {

	return db.Model(User{}).Where("name=?", userName).Update("avatar_url", avatar).Error
}
func AuthUpdate(userName, id_card, real_name, id_address string) error {

	return db.Model(User{}).Where("name=?", userName).Updates(User{Id_card: id_card, Real_name: real_name, Id_address: id_address}).Error
}

func SetupMysql() {
	config, err := conf.LoadConfig("./conf/")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	fmt.Println(config.DB_SOURCE)
	db, err = gorm.Open(mysql.Open(config.DB_SOURCE), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		SingularTable: true, //设置此属性关闭复数表名
	}})
	if err != nil {
		log.Panicln("连接数据库失败, error=" + err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Panicln("连接数据库失败, error=" + err.Error())
		return
	}
	sqlDB.SetMaxIdleConns(10)     //连接池最大允许的空闲连接数
	sqlDB.SetMaxOpenConns(100)    //设置数据库连接池最大连接数
	sqlDB.SetConnMaxLifetime(200) //设置数据库连接池可重用链接得最大时间长度
	fmt.Println("Mysql链接成功！")
}
