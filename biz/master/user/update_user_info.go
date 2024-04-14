package user

import (
	"context"

	"github.com/EquaApps/frp/biz/master/client"
	"github.com/EquaApps/frp/common"
	"github.com/EquaApps/frp/dao"
	"github.com/EquaApps/frp/models"
	"github.com/EquaApps/frp/pb"
	"github.com/EquaApps/frp/utils"
	"github.com/sirupsen/logrus"
)

func UpdateUserInfoHander(c context.Context, req *pb.UpdateUserInfoRequest) (*pb.UpdateUserInfoResponse, error) {
	var (
		userInfo = common.GetUserInfo(c)
	)

	if !userInfo.Valid() {
		return &pb.UpdateUserInfoResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid user"},
		}, nil
	}
	newUserEntity := userInfo.(*models.UserEntity)
	newUserInfo := req.GetUserInfo()

	if newUserInfo.GetEmail() != "" {
		newUserEntity.Email = newUserInfo.GetEmail()
	}

	if newUserInfo.GetRawPassword() != "" {
		hashedPassword, err := utils.HashPassword(newUserInfo.GetRawPassword())
		if err != nil {
			logrus.WithError(err).Errorf("cannot hash password")
			return nil, err
		}
		newUserEntity.Password = hashedPassword
	}

	if newUserInfo.GetUserName() != "" {
		newUserEntity.UserName = newUserInfo.GetUserName()
	}

	if newUserInfo.GetToken() != "" {
		newUserEntity.Token = newUserInfo.GetToken()
	}

	if err := dao.UpdateUser(userInfo, newUserEntity); err != nil {
		return &pb.UpdateUserInfoResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()},
		}, err
	}

	go func() {
		newUser, err := dao.GetUserByUserID(userInfo.GetUserID())
		if err != nil {
			logrus.WithError(err).Errorf("cannot get user")
			return
		}

		if err := client.SyncTunnel(c, newUser); err != nil {
			logrus.WithError(err).Errorf("cannot sync tunnel, user need to retry update")
		}
	}()

	return &pb.UpdateUserInfoResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
