package server

import (
	service_user_register "codetube.cn/proto/service-user-register"
)

type UserRegisterServer struct {
	service_user_register.UnimplementedUserRegisterServer
}

func NewUserRegisterServer() *UserRegisterServer {
	return &UserRegisterServer{}
}
