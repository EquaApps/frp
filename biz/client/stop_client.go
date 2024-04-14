package client

import (
	"context"

	"github.com/johncoker233/frpaaa/pb"
	"github.com/johncoker233/frpaaa/tunnel"
	"github.com/sirupsen/logrus"
)

func StopFRPCHandler(ctx context.Context, req *pb.StopFRPCRequest) (*pb.StopFRPCResponse, error) {
	logrus.Infof("client get a stop client request, origin is: [%+v]", req)

	tunnel.GetClientController().Stop(req.GetClientId())

	return &pb.StopFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
