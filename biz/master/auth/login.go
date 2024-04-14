package auth

import (
	"context"
	"fmt"

	"github.com/EquaApps/frp/conf"
	"github.com/EquaApps/frp/dao"
	"github.com/EquaApps/frp/middleware"
	"github.com/EquaApps/frp/pb"
	"github.com/gin-gonic/gin"
)

func LoginHandler(c context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()
	ok, user, err := dao.CheckUserPassword(username, password)
	if err != nil {
		return nil, err
	}

	if !ok {
		return &pb.LoginResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid username or password"},
		}, nil
	}

	tokenStr := conf.GetCommonJWT(fmt.Sprint(user.GetUserID()))

	ginCtx := c.(*gin.Context)
	middleware.PushTokenStr(ginCtx, tokenStr)

	return &pb.LoginResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		Token:  &tokenStr,
	}, nil
}
