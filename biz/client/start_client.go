package client

import (
	"context"

	"github.com/johncoker233/frpaaa/pb"
	"github.com/johncoker233/frpaaa/tunnel"
	"github.com/sirupsen/logrus"
)

func StartFRPCHandler(ctx context.Context, req *pb.StartFRPCRequest) (*pb.StartFRPCResponse, error) {
	logrus.Infof("client get a start client request, origin is: [%+v]", req)

	tunnel.GetClientController().Run(req.GetClientId())

	return &pb.StartFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
