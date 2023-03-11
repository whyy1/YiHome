package handler

import (
	"context"
	"fmt"
	"strconv"
	"userorder/model"
	pb "userorder/proto"
	"userorder/utils"
)

type Userorder struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Userorder) CreateOrder(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	//获取到相关数据,插入到数据库
	orderId, err := model.InsertOrder(req.HouseId, req.StartDate, req.EndDate, req.UserName)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	var orderData pb.OrderData
	orderData.OrderId = strconv.Itoa(orderId)

	rsp.Data = &orderData

	return nil
}

func (e *Userorder) GetOrderInfo(ctx context.Context, req *pb.GetReq, resp *pb.GetResp) error {
	//要根据传入数据获取订单信息   mysql
	respData, err := model.GetOrderInfo(req.UserName, req.Role)
	if err != nil {
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	var getData pb.GetData
	getData.Orders = respData
	resp.Data = &getData

	return nil
}

func (e *Userorder) UpdateStatus(ctx context.Context, req *pb.UpdateReq, resp *pb.UpdateResp) error {
	//根据传入数据,更新订单状态
	err := model.UpdateStatus(req.Action, req.Id, req.Reason)
	if err != nil {
		fmt.Println("更新订单装填错误", err)
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	return nil
}

func (e *Userorder) UpdateComment(ctx context.Context, req *pb.CommentReq, resp *pb.CommentResp) error {
	fmt.Println("接收到的参数", req)

	//根据传入数据,更新订单评论
	err := model.UpdateComment(req.OrderId, req.Comment)
	if err != nil {
		fmt.Println("更新订单装填错误", err)
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	return nil
}
