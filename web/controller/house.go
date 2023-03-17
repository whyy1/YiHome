package controller

import (
	houseMicro "YiHome/proto/house"
	"YiHome/utils"
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/storage"
	uid "github.com/satori/go.uuid"
	"go-micro.dev/v4/client"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

func GetUserHouses(ctx *gin.Context) {
	//获取用户名
	userName := sessions.Default(ctx).Get("userName")
	microClient := houseMicro.NewHouseService("house", utils.NewMicroClient().Client())
	//调用远程服务
	resp, _ := microClient.GetHouseInfo(context.TODO(), &houseMicro.GetReq{UserName: userName.(string)})
	fmt.Println("接收到的参数", resp)

	ctx.JSON(http.StatusOK, resp)
}

func PostHouses(ctx *gin.Context) {
	var house HouseStu

	if err := ctx.Bind(&house); err != nil {
		log.Println("loginData参数解析错误，err为：", err)
		ctx.JSON(http.StatusInternalServerError, fmt.Sprint("参数解析错误，", err))
		return
	}

	userName := sessions.Default(ctx).Get("userName")
	var opts client.CallOption = func(o *client.CallOptions) {
		o.RequestTimeout = time.Second * 30
		o.DialTimeout = time.Second * 30
	}
	microClient := houseMicro.NewHouseService("house", utils.NewMicroClient().Client())
	resp, err := microClient.PubHouse(context.TODO(), &houseMicro.Request{
		Acreage:   house.Acreage,
		Address:   house.Address,
		AreaId:    house.AreaId,
		Beds:      house.Beds,
		Capacity:  house.Capacity,
		Deposit:   house.Deposit,
		Facility:  house.Facility,
		MaxDays:   house.MaxDays,
		MinDays:   house.MinDays,
		Price:     house.Price,
		RoomCount: house.RoomCount,
		Title:     house.Title,
		Unit:      house.Unit,
		UserName:  userName.(string),
	}, opts)
	fmt.Println("resp==", resp, err)
	//返回数据
	ctx.JSON(http.StatusOK, resp)
}
func PostHousesImage(ctx *gin.Context) {
	houseId := ctx.Param("id")

	fileHeader, err := ctx.FormFile("house_image")
	if err != nil {
		log.Println("上传图片出错，错误为：", err)
		ctx.JSON(http.StatusInternalServerError, "上传图片出错")
	}

	filename := uid.NewV4()
	fmt.Println("生成的uuid为", filename)
	path := fmt.Sprintf("YiHome/houseImage/%v", filename)
	savepath := path + filepath.Ext(fileHeader.Filename)
	token := utils.NewUpToken(path, savepath)

	if err := utils.PutFile(token, path, fileHeader); err != nil {
		log.Fatal("文件上传失败", err)
		ctx.JSON(http.StatusInternalServerError, "使用七牛云上传错误")
	}
	imgUrl := storage.MakePublicURL("http://cdn.whyy1.top", savepath)
	microClient := houseMicro.NewHouseService("house", utils.NewMicroClient().Client())
	resp, err := microClient.UploadHouseImg(context.TODO(), &houseMicro.ImgReq{
		HouseId: houseId,
		ImgUrl:  imgUrl,
	})
	if err != nil {
		log.Println("客户端调用house微服务UploadHouseImg函数失败，err为：", err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	response := Resp{
		Errno:  resp.Errno,
		Errmsg: resp.Errmsg,
		Data: map[string]string{
			"url": imgUrl,
		},
	}
	ctx.JSON(http.StatusOK, response)

}
func GetHousesInfo(ctx *gin.Context) {
	houseId := ctx.Param("id")
	//校验数据
	if houseId == "" {
		fmt.Println("获取数据错误")
		return
	}
	userName := sessions.Default(ctx).Get("userName")
	//处理数据
	microClient := houseMicro.NewHouseService("house", utils.NewMicroClient().Client())
	//调用远程服务
	resp, _ := microClient.GetHouseDetail(context.TODO(), &houseMicro.DetailReq{
		HouseId:  houseId,
		UserName: userName.(string),
	})

	//返回数据
	ctx.JSON(http.StatusOK, resp)

}
func PutHousesInfo(ctx *gin.Context) {
	var house HouseStu
	house_id := ctx.Param("id")
	//校验数据
	if house_id == "" {
		resp := Response{
			Errno:  utils.RECODE_DATAERR,
			Errmsg: "获取不到房屋ID",
		}
		ctx.JSON(http.StatusOK, resp)
		return
	}

	if err := ctx.Bind(&house); err != nil {
		log.Println("loginData参数解析错误，err为：", err)
		ctx.JSON(http.StatusInternalServerError, fmt.Sprint("参数解析错误，", err))
		return
	}
	fmt.Println(house)

	userName := sessions.Default(ctx).Get("userName")
	fmt.Println(userName)
	//创建微服务进行update操作
	microClient := houseMicro.NewHouseService("house", utils.NewMicroClient().Client())
	var resp, err = microClient.UpdateHouse(context.TODO(), &houseMicro.UpdateResq{
		Acreage:   house.Acreage,
		Address:   house.Address,
		AreaId:    house.AreaId,
		Beds:      house.Beds,
		Capacity:  house.Capacity,
		Deposit:   house.Deposit,
		Facility:  house.Facility,
		MaxDays:   house.MaxDays,
		MinDays:   house.MinDays,
		Price:     house.Price,
		RoomCount: house.RoomCount,
		Title:     house.Title,
		Unit:      house.Unit,
		UserName:  userName.(string),
		HouseId:   house_id,
	})

	if err != nil {
		log.Println("调用微服务House中UpdateHouse出现错误，err为：", err)
	}
	//返回数据
	ctx.JSON(http.StatusOK, resp)
}
func DeleteHousesInfo(ctx *gin.Context) {
	house_id := ctx.Param("id")
	userName := sessions.Default(ctx).Get("userName")
	fmt.Println(house_id, userName)
	microClient := houseMicro.NewHouseService("house", utils.NewMicroClient().Client())
	var resp, err = microClient.DeleteHouse(context.TODO(), &houseMicro.DeleteResq{
		HouseId:  house_id,
		UserName: userName.(string),
	})
	if err != nil {
		log.Println("调用微服务House中UpdateHouse出现错误，err为：", err)
	}
	//返回数据
	ctx.JSON(http.StatusOK, resp)
}
func GetHouses(ctx *gin.Context) {
	//获取数据
	//areaId
	aid := ctx.Query("aid")
	//start day
	sd := ctx.Query("sd")
	//end day
	ed := ctx.Query("ed")
	//排序方式
	sk := ctx.Query("sk")
	//page  第几页
	//ctx.Query("p")
	//校验数据
	if aid == "" || sd == "" || ed == "" || sk == "" {
		fmt.Println("传入数据不完整")
		return
	}

	microClient := houseMicro.NewHouseService("house", utils.NewMicroClient().Client())
	//调用远程服务
	resp, _ := microClient.SearchHouse(context.TODO(), &houseMicro.SearchReq{
		Aid: aid,
		Sd:  sd,
		Ed:  ed,
		Sk:  sk,
	})

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}
