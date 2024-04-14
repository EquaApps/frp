package auth

import (
	"context"
	"fmt"

	"github.com/johncoker233/frpaaa/conf"
	"github.com/johncoker233/frpaaa/dao"
	"github.com/johncoker233/frpaaa/models"
	"github.com/johncoker233/frpaaa/pb"
	"github.com/johncoker233/frpaaa/utils"
	"github.com/google/uuid"
)

func RegisterHandler(c context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()
	email := req.GetEmail()

	if username == "" || password == "" || email == "" {
		return &pb.RegisterResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid username or password or email"},
		}, fmt.Errorf("invalid username or password or email")
	}

	userCount, err := dao.AdminCountUsers()
	if err != nil {
		return &pb.RegisterResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()},
		}, err
	}

	if !conf.Get().App.EnableRegister && userCount > 0 {
		return &pb.RegisterResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "register is disabled"},
		}, fmt.Errorf("register is disabled")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return &pb.RegisterResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()},
		}, err
	}

	newUser := &models.UserEntity{
		UserName: username,
		Password: hashedPassword,
		Email:    email,
		Status:   models.STATUS_NORMAL,
		Role:     models.ROLE_NORMAL,
		Token:    uuid.New().String(),
	}

	err = dao.CreateUser(newUser)
	if err != nil {
		return &pb.RegisterResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()},
		}, err
	}

	return &pb.RegisterResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
