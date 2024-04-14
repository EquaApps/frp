package server

import (
	"context"

	"github.com/EquaApps/frp/common"
	"github.com/EquaApps/frp/dao"
	"github.com/EquaApps/frp/models"
	"github.com/EquaApps/frp/pb"
	"github.com/samber/lo"
)

func ListServersHandler(c context.Context, req *pb.ListServersRequest) (*pb.ListServersResponse, error) {
	var (
		userInfo = common.GetUserInfo(c)
		page     = int(req.GetPage())
		pageSize = int(req.GetPageSize())
	)

	if !userInfo.Valid() {
		return &pb.ListServersResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid user"},
		}, nil
	}

	servers, err := dao.ListServers(userInfo, page, pageSize)
	if err != nil {
		return nil, err
	}

	serverCounts, err := dao.CountServers(userInfo)
	if err != nil {
		return nil, err
	}

	return &pb.ListServersResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		Servers: lo.Map(servers, func(c *models.ServerEntity, _ int) *pb.Server {
			return &pb.Server{
				Id:      lo.ToPtr(c.ServerID),
				Config:  lo.ToPtr(string(c.ConfigContent)),
				Secret:  lo.ToPtr(c.ConnectSecret),
				Ip:      lo.ToPtr(c.ServerIP),
				Comment: lo.ToPtr(c.Comment),
			}
		}),
		Total: lo.ToPtr(int32(serverCounts)),
	}, nil
}
