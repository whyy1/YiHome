package handler

import (
	"context"
	"fmt"
	"house/model"
	pb "house/proto"
	"house/utils"
	"log"
	"strconv"
)

type House struct{}

func (e *House) PubHouse(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	//上传房屋业务  把获取到的房屋数据插入数据库
	houseId, err := model.AddHouse(req)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	var h pb.HouseData
	h.HouseId = strconv.Itoa(houseId)
	rsp.Data = &h
	fmt.Println("上传房屋调用成功，resp=:", rsp)
	return nil
}
func (e *House) UploadHouseImg(ctx context.Context, req *pb.ImgReq, resp *pb.ImgResp) error {
	err := model.SaveHouseImg(req.HouseId, req.ImgUrl)
	if err != nil {
		log.Println("保存房屋图片出错，SaveHouseImg函数，err为:", err)
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	return nil
}
func (e *House) GetHouseInfo(ctx context.Context, req *pb.GetReq, resp *pb.GetResp) error {
	//根据用户名获取所有的房屋数据
	houseInfos, err := model.GetUserHouse(req.UserName)
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	var getData pb.GetData
	getData.Houses = houseInfos

	resp.Data = &getData
	fmt.Println("房源信息", resp.Data)
	return nil
}

func (e *House) GetHouseDetail(ctx context.Context, req *pb.DetailReq, resp *pb.DetailResp) error {
	//根据houseId获取所有的返回数据
	respData, err := model.GetHouseDetail(req.HouseId, req.UserName)
	if err != nil {
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	resp.Data = &respData

	return nil
}

func (e *House) GetIndexHouse(ctx context.Context, req *pb.IndexReq, resp *pb.GetResp) error {
	//获取房屋信息
	houseResp, err := model.GetIndexHouse()
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	resp.Data = &pb.GetData{Houses: houseResp}

	return nil
}

func (e *House) SearchHouse(ctx context.Context, req *pb.SearchReq, resp *pb.GetResp) error {
	//根据传入的参数,查询符合条件的房屋信息
	houseResp, err := model.SearchHouse(req.Aid, req.Sd, req.Ed, req.Sk)
	if err != nil {
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	resp.Data = &pb.GetData{Houses: houseResp}
	return nil
}

func (e *House) UpdateHouse(ctx context.Context, req *pb.UpdateResq, resp *pb.UpdateResp) error {
	//修改房源信息
	if err := model.UpdateHouse(req); err != nil {
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText("数据库修改数据时错误")
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	return nil
}

func (e *House) DeleteHouse(ctx context.Context, req *pb.DeleteResq, resp *pb.DeleteResp) error {
	//修改房源信息
	if err := model.DeleteHouse(req); err != nil {
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText("数据库修改数据时错误")
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	return nil
}
