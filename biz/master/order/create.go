package order

import (
    "context"
    "github.com/EquaApps/frp/common"
    "github.com/EquaApps/frp/dao"
    "github.com/EquaApps/frp/models"
    "github.com/EquaApps/frp/pb"
)




func CreateOrderHandler(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return &pb.CreateOrderResponse{
			Status: &pb.Status{
				Code:    pb.RespCode_RESP_CODE_INVALID,
				Message: "invalid user",
			},
		}, nil
	}

	bandwidth := req.GetBandwidth()
	if bandwidth <= 0 {
		return &pb.CreateOrderResponse{
			Status: &pb.Status{
				Code:    pb.RespCode_RESP_CODE_INVALID,
				Message: "invalid bandwidth",
			},
		}, nil
	}

	comment := req.GetComment()

	order, err := dao.CreateOrder(userInfo.GetUserID(), bandwidth, userInfo.GetUserID(), comment)
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{
		Status: &pb.Status{
			Code:    pb.RespCode_RESP_CODE_SUCCESS,
			Message: "ok",
		},
		OrderId: order.OrderID,
	}, nil
}
