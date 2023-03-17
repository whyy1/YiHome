package model

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	"time"
	"userorder/conf"
	pb "userorder/proto"
)

var (
	GlobalRedis *redis.Pool
	GlobalDB    *gorm.DB
)

type UserData struct {
	Id int
}

func InsertOrder(houseId, beginDate, endDate, userName string) (int, error) {
	//获取插入对象
	var order OrderHouse

	//给对象赋值
	hid, _ := strconv.Atoi(houseId)
	order.HouseId = uint(hid)

	//把string类型的时间转换为time类型
	bDate, _ := time.Parse("2006-01-02", beginDate)
	order.Begin_date = bDate

	eDate, _ := time.Parse("2006-01-02", endDate)
	order.End_date = eDate

	//需要userId
	/*var user User
	GlobalDB.Where("name = ?",userName).Find(&user)*/
	//select id form user where name = userName

	var userData UserData
	if err := GlobalDB.Raw("select id from user where name = ?", userName).Scan(&userData).Error; err != nil {
		fmt.Println("获取用户数据错误", err)
		return 0, err
	}

	//获取days
	dur := eDate.Sub(bDate)
	order.Days = int(dur.Hours()) / 24
	order.Status = "WAIT_ACCEPT"

	//房屋的单价和总价
	var house House
	GlobalDB.Where("id = ?", hid).Find(&house).Select("price")
	order.House_price = house.Price
	order.Amount = house.Price * order.Days

	order.UserId = uint(userData.Id)
	if err := GlobalDB.Create(&order).Error; err != nil {
		fmt.Println("插入订单失败", err)
		return 0, err
	}
	return int(order.ID), nil
}

//获取房东订单如何实现?
func GetOrderInfo(userName, role string) ([]*pb.OrdersData, error) {
	//最终需要的数据
	var orderResp []*pb.OrdersData
	//获取当前用户的所有订单
	var orders []OrderHouse

	var userData UserData
	//用原生查询的时候,查询的字段必须跟数据库中的字段保持一直
	err := GlobalDB.Debug().Raw("select id from user where name = ?", userName).Scan(&userData).Error
	if err != nil {
		fmt.Println("查询失败，err为", err)
	}

	//查询租户的所有的订单
	if role == "custom" {
		if err := GlobalDB.Where("user_id = ?", userData.Id).Find(&orders).Error; err != nil {
			fmt.Println("获取当前用户所有订单失败")
			return nil, err
		}
	} else {
		//查询房东的订单  以房东视角来查看订单
		var houses []House
		GlobalDB.Debug().Where("user_id = ?", userData.Id).Find(&houses)

		for _, v := range houses {
			var tempOrders []OrderHouse
			GlobalDB.Debug().Model(&v).Related(&tempOrders)

			orders = append(orders, tempOrders...)
		}
	}

	//循环遍历一下orders
	for _, v := range orders {
		var orderTemp pb.OrdersData
		orderTemp.OrderId = int32(v.ID)
		orderTemp.EndDate = v.End_date.Format("2006-01-02")
		orderTemp.StartDate = v.Begin_date.Format("2006-01-02")
		orderTemp.Ctime = v.CreatedAt.Format("2006-01-02")
		orderTemp.Amount = int32(v.Amount)
		orderTemp.Comment = v.Comment
		orderTemp.Days = int32(v.Days)
		orderTemp.Status = v.Status

		//关联house表
		var house House
		GlobalDB.Model(&v).Related(&house).Select("index_image_url", "title")
		orderTemp.ImgUrl = house.Index_image_url
		orderTemp.Title = house.Title

		orderResp = append(orderResp, &orderTemp)
	}
	return orderResp, nil
}

//更新订单状态
func UpdateStatus(action, id, reason string) error {
	db := GlobalDB.Model(new(OrderHouse)).Where("id = ?", id)

	if action == "accept" {
		//标示房东同意订单
		return db.Update("status", "WAIT_COMMENT").Error
	} else {
		//表示房东不同意订单  如果拒单把拒绝的原因写到comment中
		return db.Updates(map[string]interface{}{"status": "REJECTED", "comment": reason}).Error
	}
}
func UpdateComment(id, comment string) error {

	db := GlobalDB.Debug().Model(new(OrderHouse)).Where("id = ?", id)

	return db.Debug().Updates(map[string]interface{}{"status": "COMPLETE", "comment": comment}).Error
}

func InitMysql() {
	config, err := conf.LoadConfig("./conf/")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
		return
	}

	db, err := gorm.Open("mysql", config.DB_SOURCE)
	if err != nil {
		log.Fatal("Gorm打开数据库连接失败: ", err)
		panic(err)
		return
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(50)        //连接池最大允许的空闲连接数
	db.DB().SetMaxOpenConns(70)        //设置数据库连接池最大连接数
	db.DB().SetConnMaxLifetime(60 * 5) //设置数据库连接池可重用链接得最大时间长度

	GlobalDB = db
	log.Println("Mysql链接成功！")
}

func InitRedis() {
	config, err := conf.LoadConfig("./conf/")
	if err != nil {
		log.Println("初始化Redis失败，读取Redis配置文件失败: ", err)
		return
	}

	GlobalRedis = &redis.Pool{ //实例化一个连接池
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
