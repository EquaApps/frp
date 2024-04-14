package client

import (
	"context"

	"github.com/EquaApps/frp/common"
	"github.com/EquaApps/frp/dao"
	"github.com/EquaApps/frp/models"
	"github.com/EquaApps/frp/pb"
	"github.com/google/uuid"
)

func InitClientHandler(c context.Context, req *pb.InitClientRequest) (*pb.InitClientResponse, error) {
	userClientID := req.GetClientId()
	userInfo := common.GetUserInfo(c)

	if !userInfo.Valid() {
		return &pb.InitClientResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid user"},
		}, nil
	}

	if len(userClientID) == 0 {
		return &pb.InitClientResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid client id"},
		}, nil
	}

	globalClientID := common.GlobalClientID(userInfo.GetUserName(), "c", userClientID)

	if err := dao.CreateClient(userInfo,
		&models.ClientEntity{
			ClientID:      globalClientID,
			TenantID:      userInfo.GetTenantID(),
			UserID:        userInfo.GetUserID(),
			ConnectSecret: uuid.New().String(),
		}); err != nil {
		return nil, err
	}

	return &pb.InitClientResponse{
		Status:   &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		ClientId: &globalClientID,
	}, nil
}
