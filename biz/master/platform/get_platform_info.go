package platform

import (
	"github.com/johncoker233/frpaaa/common"
	"github.com/johncoker233/frpaaa/conf"
	"github.com/johncoker233/frpaaa/dao"
	"github.com/johncoker233/frpaaa/pb"
	"github.com/gin-gonic/gin"
)

func GetPlatformInfo(c *gin.Context) {
	resp, err := getPlatformInfo(c)
	if err != nil {
		common.ErrResp(c, &pb.CommonResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()}}, err.Error())
		return
	}
	common.OKResp(c, resp)
}

func getPlatformInfo(c *gin.Context) (*pb.GetPlatformInfoResponse, error) {
	userInfo := common.GetUserInfo(c)
	if !userInfo.Valid() {
		return &pb.GetPlatformInfoResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid user"},
		}, nil
	}
	totalServers, err := dao.CountServers(userInfo)
	if err != nil {
		return nil, err
	}
	totalClients, err := dao.CountClients(userInfo)
	if err != nil {
		return nil, err
	}

	unconfiguredServers, err := dao.CountUnconfiguredServers(userInfo)
	if err != nil {
		return nil, err
	}

	unconfiguredClients, err := dao.CountUnconfiguredClients(userInfo)
	if err != nil {
		return nil, err
	}

	configuredServers, err := dao.CountConfiguredServers(userInfo)
	if err != nil {
		return nil, err
	}
	configuredClients, err := dao.CountConfiguredClients(userInfo)
	if err != nil {
		return nil, err
	}
	return &pb.GetPlatformInfoResponse{
		Status:                  &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		TotalClientCount:        int32(totalClients),
		TotalServerCount:        int32(totalServers),
		UnconfiguredClientCount: int32(unconfiguredClients),
		UnconfiguredServerCount: int32(unconfiguredServers),
		ConfiguredClientCount:   int32(configuredClients),
		ConfiguredServerCount:   int32(configuredServers),
		GlobalSecret:            conf.MasterDefaultSalt(),
		MasterRpcHost:           conf.Get().Master.RPCHost,
		MasterRpcPort:           int32(conf.Get().Master.RPCPort),
		MasterApiPort:           int32(conf.Get().Master.APIPort),
		MasterApiScheme:         conf.Get().Master.APIScheme,
	}, nil
}
