package server

import (
	"codetube.cn/core/codes"
	service_user_register "codetube.cn/proto/service-user-register"
	"codetube.cn/service-user-register/components"
	"codetube.cn/service-user-register/libraries/check"
	"codetube.cn/service-user-register/libraries/password"
	"codetube.cn/service-user-register/models"
	"context"
	"log"
	"strconv"
)

// UserPassword 使用用户名和密码注册
func (s *UserRegisterServer) UserPassword(c context.Context, request *service_user_register.UsernamePasswordRequest) (*service_user_register.RegisterResultResponse, error) {
	status := codes.Success
	message := "success"

	username := request.GetUsername()
	passwd := request.GetPassword()
	id := int64(0)
	//检查用户名和密码格式
	if len(username) < 5 || len(username) > 20 {
		status = codes.UserAccountInvalid
	}
	if len(passwd) < 8 || len(passwd) > 40 {
		status = codes.UserPasswordInvalid
	}
	if status != codes.Success {
		username = ""
	}
	//检查用户名是否重复
	userExist := check.UserExistByUsername(username)
	if userExist {
		status = codes.UserAccountIsExist
	} else {
		//写入数据
		//定义一个用户，并初始化数据
		user := models.User{
			Username: username,
			Password: password.HashPassword(passwd),
			Enabled:  1,
		}
		//插入记录
		result := components.UserDB.Create(&user)

		if result.Error != nil {
			status = codes.UserRegisterInsertDbFailed
			log.Println("[err:"+strconv.Itoa(codes.UserRegisterInsertDbFailed)+"]写入用户信息失败：", result.Error)
		} else if user.ID < 1 {
			status = codes.UserRegisterFailed
		} else {
			id = int64(user.ID)
			//@todo 其他逻辑，例如广播通知其他服务
		}
	}
	return &service_user_register.RegisterResultResponse{
		Status:  int64(status),
		Message: message,
		Id:      id,
	}, nil
}
