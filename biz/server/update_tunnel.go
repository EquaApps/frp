package server

import (
	"context"
	"reflect"

	"github.com/johncoker233/frpaaa/pb"
	"github.com/johncoker233/frpaaa/services/server"
	"github.com/johncoker233/frpaaa/tunnel"
	"github.com/johncoker233/frpaaa/utils"
	"github.com/sirupsen/logrus"
)

func UpdateFrpsHander(ctx context.Context, req *pb.UpdateFRPSRequest) (*pb.UpdateFRPSResponse, error) {
	logrus.Infof("update frps, req: [%+v]", req)

	content := req.GetConfig()

	s, err := utils.LoadServerConfig(content, true)
	if err != nil {
		logrus.WithError(err).Errorf("cannot load config")
		return &pb.UpdateFRPSResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()},
		}, err
	}

	serverID := req.GetServerId()
	if cli := tunnel.GetServerController().Get(serverID); cli != nil {
		if !reflect.DeepEqual(cli.GetCommonCfg(), s) {
			cli.Stop()
			tunnel.GetClientController().Delete(serverID)
			logrus.Infof("server %s config changed, will recreate it", serverID)
		} else {
			logrus.Infof("server %s config not changed", serverID)
			return &pb.UpdateFRPSResponse{
				Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
			}, nil
		}
	}
	tunnel.GetServerController().Add(serverID, server.NewServerHandler(s))
	tunnel.GetServerController().Run(serverID)

	return &pb.UpdateFRPSResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
