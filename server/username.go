package server

import (
	"codetube.cn/core/codes"
	"codetube.cn/proto/service_user_register"
	"codetube.cn/service-user-register/components"
	"codetube.cn/service-user-register/libraries/check"
	"codetube.cn/service-user-register/libraries/password"
	"codetube.cn/service-user-register/models"
	"context"
	"github.com/google/uuid"
	"log"
	"strconv"
)

// Username 使用用户名和密码注册
func (s *UserRegisterServer) Username(c context.Context, request *service_user_register.RegisterUsernameRequest) (*service_user_register.RegisterResultResponse, error) {
	message := "success"
	status := s.checkUsernameParams(request)
	if status != codes.Success {
		return &service_user_register.RegisterResultResponse{
			Status:  int64(status),
			Message: message, //@todo 数字转文字
			Id:      "",
		}, nil
	}

	var id string
	username := request.GetUsername()
	passwd := request.GetPassword()
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
		} else if user.ID == uuid.Nil {
			status = codes.UserRegisterFailed
		} else {
			id = user.ID.String()
			//@todo 其他逻辑，例如广播通知其他服务
		}
	}
	return &service_user_register.RegisterResultResponse{
		Status:  int64(status),
		Message: message,
		Id:      id,
	}, nil
}

func (s *UserRegisterServer) checkUsernameParams(request *service_user_register.RegisterUsernameRequest) int {
	username := request.GetUsername()
	passwd := request.GetPassword()

	//检查用户名和密码格式
	if len(username) < 5 || len(username) > 20 {
		return codes.UserAccountInvalid
	}
	if !password.CheckPassword(passwd) {
		return codes.UserPasswordInvalid
	}
	return codes.Success
}
