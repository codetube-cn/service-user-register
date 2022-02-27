package server

import (
	"codetube.cn/core/codes"
	service_user_register "codetube.cn/proto/service-user-register"
	"context"
)

type UserRegisterServer struct {
	service_user_register.UnimplementedUserRegisterServer
}

func NewUserRegisterServer() *UserRegisterServer {
	return &UserRegisterServer{}
}

func (s *UserRegisterServer) UserPassword(c context.Context, request *service_user_register.UsernamePasswordRequest) (*service_user_register.RegisterResultResponse, error) {
	return &service_user_register.RegisterResultResponse{
		Status:   codes.Success,
		Message:  "success",
		Username: request.GetUsername(),
	}, nil
}
