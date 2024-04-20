package order

import (
    "context"
    "github.com/johncoker233/frpaaa/common"
    "github.com/johncoker233/frpaaa/dao"
    "github.com/johncoker233/frpaaa/models"
    "github.com/johncoker233/frpaaa/pb"
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
