package server

import (
	"context"

	"github.com/EquaApps/frp/common"
	"github.com/EquaApps/frp/dao"
	"github.com/EquaApps/frp/models"
	"github.com/EquaApps/frp/pb"
	"github.com/google/uuid"
)

func InitServerHandler(c context.Context, req *pb.InitServerRequest) (*pb.InitServerResponse, error) {
	var (
		userServerID = req.GetServerId()
		serverIP     = req.GetServerIp()
		userInfo     = common.GetUserInfo(c)
	)

	if !userInfo.Valid() {
		return &pb.InitServerResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid user"},
		}, nil
	}

	if len(userServerID) == 0 || len(serverIP) == 0 {
		return &pb.InitServerResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "request invalid"},
		}, nil
	}

	globalServerID := common.GlobalClientID(userInfo.GetUserName(), "s", userServerID)

	if err := dao.CreateServer(userInfo,
		&models.ServerEntity{
			ServerID:      globalServerID,
			TenantID:      userInfo.GetTenantID(),
			UserID:        userInfo.GetUserID(),
			ConnectSecret: uuid.New().String(),
			ServerIP:      serverIP,
		}); err != nil {
		return nil, err
	}

	return &pb.InitServerResponse{
		Status:   &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		ServerId: &globalServerID,
	}, nil
}
