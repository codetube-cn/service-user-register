package server

import (
	"codetube.cn/proto/service_user_register"
)

type UserRegisterServer struct {
	service_user_register.UnimplementedUserRegisterServer
}

func NewUserRegisterServer() *UserRegisterServer {
	return &UserRegisterServer{}
}
