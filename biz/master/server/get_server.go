package server

import (
	"context"

	"github.com/johncoker233/frpaaa/common"
	"github.com/johncoker233/frpaaa/dao"
	"github.com/johncoker233/frpaaa/pb"
	"github.com/samber/lo"
)

func GetServerHandler(c context.Context, req *pb.GetServerRequest) (*pb.GetServerResponse, error) {
	var (
		userServerID = req.GetServerId()
		userInfo     = common.GetUserInfo(c)
	)

	if !userInfo.Valid() {
		return &pb.GetServerResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid user"},
		}, nil
	}

	if len(userServerID) == 0 {
		return &pb.GetServerResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid client id"},
		}, nil
	}

	serverEntity, err := dao.GetServerByServerID(userInfo, userServerID)
	if err != nil {
		return nil, err
	}

	return &pb.GetServerResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		Server: &pb.Server{
			Id:      lo.ToPtr(serverEntity.ServerID),
			Config:  lo.ToPtr(string(serverEntity.ConfigContent)),
			Secret:  lo.ToPtr(serverEntity.ConnectSecret),
			Comment: lo.ToPtr(serverEntity.Comment),
			Ip:      lo.ToPtr(serverEntity.ServerIP),
		},
	}, nil
}
