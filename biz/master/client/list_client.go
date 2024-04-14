package client

import (
	"context"

	"github.com/johncoker233/frpaaa/common"
	"github.com/johncoker233/frpaaa/dao"
	"github.com/johncoker233/frpaaa/models"
	"github.com/johncoker233/frpaaa/pb"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func ListClientsHandler(ctx context.Context, req *pb.ListClientsRequest) (*pb.ListClientsResponse, error) {
	logrus.Infof("list client, req: [%+v]", req)

	var (
		userInfo = common.GetUserInfo(ctx)
	)

	if !userInfo.Valid() {
		return &pb.ListClientsResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid user"},
		}, nil
	}

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	clients, err := dao.ListClients(userInfo, page, pageSize)
	if err != nil {
		return nil, err
	}

	clientCounts, err := dao.CountClients(userInfo)
	if err != nil {
		return nil, err
	}

	respClients := lo.Map(clients, func(c *models.ClientEntity, _ int) *pb.Client {
		return &pb.Client{
			Id:       lo.ToPtr(c.ClientID),
			Secret:   lo.ToPtr(c.ConnectSecret),
			Config:   lo.ToPtr(string(c.ConfigContent)),
			ServerId: lo.ToPtr(c.ServerID),
			Stopped:  lo.ToPtr(c.Stopped),
			Comment:  lo.ToPtr(c.Comment),
		}
	})

	return &pb.ListClientsResponse{
		Status:  &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		Clients: respClients,
		Total:   lo.ToPtr(int32(clientCounts)),
	}, nil
}
