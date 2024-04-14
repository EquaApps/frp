package server

import (
	"context"

	"github.com/EquaApps/frp/common"
	"github.com/EquaApps/frp/dao"
	"github.com/EquaApps/frp/pb"
)

func DeleteServerHandler(c context.Context, req *pb.DeleteServerRequest) (*pb.DeleteServerResponse, error) {
	var (
		userServerID = req.GetServerId()
		userInfo     = common.GetUserInfo(c)
	)

	if !userInfo.Valid() {
		return &pb.DeleteServerResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid user"},
		}, nil
	}

	if len(userServerID) == 0 {
		return &pb.DeleteServerResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid client id"},
		}, nil
	}

	if err := dao.DeleteServer(userInfo, userServerID); err != nil {
		return nil, err
	}

	return &pb.DeleteServerResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
