package model

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"house/conf"
	house "house/proto"
	"log"
	"strconv"
	"time"
)

var (
	GlobalRedis *redis.Pool
	GlobalDB    *gorm.DB
)

func AddHouse(request *house.Request) (int, error) {
	var houseInfo House
	//给house赋值
	houseInfo.Address = request.Address

	//根据userName获取userId
	var user User
	if err := GlobalDB.Debug().Where("name = ?", request.UserName).Find(&user).Error; err != nil {
		fmt.Println("查询当前用户失败", err)
		return 0, err
	}

	//sql中一对多插入,只是给外键赋值
	houseInfo.UserId = uint(user.ID)
	houseInfo.Title = request.Title
	//类型转换
	price, _ := strconv.Atoi(request.Price)
	roomCount, _ := strconv.Atoi(request.RoomCount)
	houseInfo.Price = price
	houseInfo.Room_count = roomCount
	houseInfo.Unit = request.Unit
	houseInfo.Capacity, _ = strconv.Atoi(request.Capacity)
	houseInfo.Beds = request.Beds
	houseInfo.Deposit, _ = strconv.Atoi(request.Deposit)
	houseInfo.Min_days, _ = strconv.Atoi(request.MinDays)
	houseInfo.Max_days, _ = strconv.Atoi(request.MaxDays)
	houseInfo.Acreage, _ = strconv.Atoi(request.MaxDays)
	//一对多插入
	areaId, _ := strconv.Atoi(request.AreaId)
	houseInfo.AreaId = uint(areaId)

	facilityMap := make(map[int]Facility)
	var facility []Facility
	GlobalDB.Debug().Find(&facility)
	fmt.Println("获取到的所有家具", facility)
	for _, v := range facility {
		facilityMap[v.Id] = v
	}
	fmt.Println("map中的值", facilityMap)

	//request.Facility    所有的家具  房屋
	for _, v := range request.Facility {
		id, _ := strconv.Atoi(v)

		v, ok := facilityMap[id]
		if ok {
			houseInfo.Facilities = append(houseInfo.Facilities, &v)
		} else {
			fmt.Println("查询不到家具id")
		}
		//查询到了数据
		//houseInfo.Facilities = append(houseInfo.Facilities, &fac)
	}

	fmt.Println("houseInfo数据为：", houseInfo)
	if err := GlobalDB.Debug().Create(&houseInfo).Error; err != nil {
		log.Println("插入房屋信息失败", err)
		return 0, err
	}
	fmt.Println("插入房屋信息成功")
	return int(houseInfo.ID), nil
}
func SaveHouseImg(houseId, imgUrl string) error {
	var houseInfo House
	if err := GlobalDB.First(&houseInfo, houseId).Error; err != nil {
		log.Println("保存房屋图片时，查询不到房屋信息", err)
		return err
	}

	if houseInfo.Index_image_url == "" {
		//说明没有上传过图片  现在上传的图片是主图
		return GlobalDB.Model(new(House)).Where("id = ?", houseId).
			Update("index_image_url", imgUrl).Error
	}
	hId, _ := strconv.Atoi(houseId)
	houseImg := HouseImage{
		Url:     imgUrl,
		HouseId: uint(hId),
	}

	return GlobalDB.Create(&houseImg).Error
}
func GetUserHouse(userName string) ([]*house.Houses, error) {
	var houseInfos []*house.Houses

	//有用户名
	var user User
	if err := GlobalDB.Where("name = ?", userName).Find(&user).Error; err != nil {
		fmt.Println("获取当前用户信息错误", err)
		return nil, err
	}

	//房源信息   一对多查询
	var houses []House
	if err := GlobalDB.Debug().Model(&user).Related(&houses).Error; err != nil {
		fmt.Println("联合查询错误错误", err)
	}

	for _, v := range houses {
		var houseInfo house.Houses
		houseInfo.Title = v.Title
		houseInfo.Address = v.Address
		houseInfo.Ctime = v.CreatedAt.Format("2006-01-02 15:04:05")
		houseInfo.HouseId = int32(v.ID)
		houseInfo.ImgUrl = v.Index_image_url
		houseInfo.OrderCount = int32(v.Order_count)
		houseInfo.Price = int32(v.Price)
		houseInfo.RoomCount = int32(v.Room_count)
		houseInfo.UserAvatar = user.Avatar_url

		//获取地域信息
		var area Area
		//related函数可以是以主表关联从表,也可以是以从表关联主表
		GlobalDB.Where("id = ?", v.AreaId).Find(&area)
		houseInfo.AreaName = area.Name

		houseInfos = append(houseInfos, &houseInfo)
	}

	return houseInfos, nil
}

//获取房屋详情
func GetHouseDetail(houseId, userName string) (house.DetailData, error) {

	var respData house.DetailData

	//给houseDetail赋值
	var houseDetail house.HouseDetail

	var houseInfo House
	if err := GlobalDB.Where("id = ?", houseId).Find(&houseInfo).Error; err != nil {
		fmt.Println("查询房屋信息错误", err)
		return respData, err
	}
	{
		houseDetail.Acreage = int32(houseInfo.Acreage)
		houseDetail.Address = houseInfo.Address
		houseDetail.Beds = houseInfo.Beds
		houseDetail.Capacity = int32(houseInfo.Capacity)
		houseDetail.Deposit = int32(houseInfo.Deposit)
		houseDetail.Hid = int32(houseInfo.ID)
		houseDetail.MaxDays = int32(houseInfo.Max_days)
		houseDetail.MinDays = int32(houseInfo.Min_days)
		houseDetail.Price = int32(houseInfo.Price)
		houseDetail.RoomCount = int32(houseInfo.Room_count)
		houseDetail.Title = houseInfo.Title
		houseDetail.Unit = houseInfo.Unit
		if houseInfo.Index_image_url != "" {
			houseDetail.ImgUrls = append(houseDetail.ImgUrls, houseInfo.Index_image_url)
		}
	}

	//评论在order表
	var orders []OrderHouse
	if err := GlobalDB.Model(&houseInfo).Related(&orders).Error; err != nil {
		fmt.Println("查询房屋评论信息", err)
		return respData, err
	}
	//var comments []*house.CommentData
	for _, v := range orders {
		var commentTemp house.CommentData
		commentTemp.Comment = v.Comment
		commentTemp.Ctime = v.CreatedAt.Format("2006-01-02 15:04:05")
		var tempUser User
		GlobalDB.Model(&v).Related(&tempUser)
		commentTemp.UserName = tempUser.Name

		houseDetail.Comments = append(houseDetail.Comments, &commentTemp)
	}

	//获取房屋的家具信息  多对多查询
	var facs []Facility
	if err := GlobalDB.Model(&houseInfo).Related(&facs, "Facilities").Error; err != nil {
		fmt.Println("查询房屋家具信息错误", err)
		return respData, err
	}
	for _, v := range facs {
		houseDetail.Facilities = append(houseDetail.Facilities, int32(v.Id))
	}

	//获取副图片  幅图找不到算不算错
	var imgs []HouseImage
	if err := GlobalDB.Model(&houseInfo).Related(&imgs).Error; err != nil {
		fmt.Println("该房屋只有主图", err)
	}

	for _, v := range imgs {
		if len(imgs) != 0 {
			houseDetail.ImgUrls = append(houseDetail.ImgUrls, v.Url)
		}
	}

	//获取房屋所有者信息
	var user User
	if err := GlobalDB.Model(&houseInfo).Related(&user).Error; err != nil {
		fmt.Println("查询房屋所有者信息错误", err)
		return respData, err
	}
	houseDetail.UserName = user.Name
	houseDetail.UserAvatar = user.Avatar_url
	houseDetail.UserId = int32(user.ID)

	respData.House = &houseDetail

	//获取当前浏览人信息
	var nowUser User
	if err := GlobalDB.Where("name = ?", userName).Find(&nowUser).Error; err != nil {
		fmt.Println("查询当前浏览人信息错误", err)
		return respData, err
	}
	respData.UserId = int32(nowUser.ID)
	return respData, nil
}
func GetIndexHouse() ([]*house.Houses, error) {

	var housesResp []*house.Houses

	var houses []House
	if err := GlobalDB.Limit(5).Find(&houses).Error; err != nil {
		fmt.Println("获取房屋信息失败", err)
		return nil, err
	}

	for _, v := range houses {
		var houseTemp house.Houses
		houseTemp.Address = v.Address
		//根据房屋信息获取地域信息
		var area Area
		var user User

		GlobalDB.Model(&v).Related(&area).Related(&user)

		houseTemp.AreaName = area.Name
		houseTemp.Ctime = v.CreatedAt.Format("2006-01-02 15:04:05")
		houseTemp.HouseId = int32(v.ID)
		houseTemp.ImgUrl = v.Index_image_url
		houseTemp.OrderCount = int32(v.Order_count)
		houseTemp.Price = int32(v.Price)
		houseTemp.RoomCount = int32(v.Room_count)
		houseTemp.Title = v.Title
		houseTemp.UserAvatar = user.Avatar_url

		housesResp = append(housesResp, &houseTemp)
	}

	return housesResp, nil
}

//搜索房屋
func SearchHouse(areaId, sd, ed, sk string) ([]*house.Houses, error) {
	var houseInfos []House

	//   minDays  <  (结束时间  -  开始时间) <  max_days
	//计算一个差值  先把string类型转为time类型
	sdTime, _ := time.Parse("2006-01-02", sd)
	edTime, _ := time.Parse("2006-01-02", ed)
	dur := edTime.Sub(sdTime)

	err := GlobalDB.Debug().Where("area_id = ?", areaId).
		Where("min_days <= ?", dur.Hours()/24).
		Where("max_days >= ? or max_days =0 ", dur.Hours()/24).
		Order("created_at desc").Find(&houseInfos).Error
	if err != nil {
		fmt.Println("搜索房屋失败", err)
		return nil, err
	}

	//获取[]*house.Houses
	var housesResp []*house.Houses

	for _, v := range houseInfos {
		var houseTemp house.Houses
		houseTemp.Address = v.Address
		//根据房屋信息获取地域信息
		var area Area
		var user User

		GlobalDB.Debug().Model(&v).Related(&area).Related(&user)

		houseTemp.AreaName = area.Name
		houseTemp.Ctime = v.CreatedAt.Format("2006-01-02 15:04:05")
		houseTemp.HouseId = int32(v.ID)
		houseTemp.ImgUrl = v.Index_image_url
		houseTemp.OrderCount = int32(v.Order_count)
		houseTemp.Price = int32(v.Price)
		houseTemp.RoomCount = int32(v.Room_count)
		houseTemp.Title = v.Title
		houseTemp.UserAvatar = user.Avatar_url

		housesResp = append(housesResp, &houseTemp)

	}

	return housesResp, nil
}

func InitMysql() {
	config, err := conf.LoadConfig("./conf/")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
		return
	}
	db, err := gorm.Open("mysql", config.DB_SOURCE)
	if err != nil {
		panic(err)
	}
	//db, err := gorm.Open("YiHome", config.DB_SOURCE)
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(50)        //连接池最大允许的空闲连接数
	db.DB().SetMaxOpenConns(70)        //设置数据库连接池最大连接数
	db.DB().SetConnMaxLifetime(60 * 5) //设置数据库连接池可重用链接得最大时间长度

	GlobalDB = db
	fmt.Println("Mysql链接成功！")
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
