package server

import (
	"codetube.cn/core/codes"
	"codetube.cn/core/libraries"
	service_user_register "codetube.cn/proto/service-user-register"
	"codetube.cn/service-user-register/components"
	"codetube.cn/service-user-register/libraries/check"
	"codetube.cn/service-user-register/libraries/password"
	"codetube.cn/service-user-register/models"
	"context"
	"log"
	"strconv"
)

// Email 使用邮箱注册
func (s *UserRegisterServer) Email(c context.Context, request *service_user_register.EmailRequest) (*service_user_register.RegisterResultResponse, error) {
	message := "success"
	status := s.checkEmailParams(request)
	if status != codes.Success {
		return &service_user_register.RegisterResultResponse{
			Status:  int64(status),
			Message: message, //@todo 数字转文字
			Id:      0,
		}, nil
	}

	email := request.GetEmail()
	passwd := request.GetPassword()
	id := int64(0)
	//检查邮箱是否重复
	userExist := check.UserExistByEmail(email)
	if userExist {
		status = codes.UserAccountIsExist
	} else {
		//写入数据
		//定义一个用户，并初始化数据
		user := models.User{
			Email:    email,
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

func (s *UserRegisterServer) checkEmailParams(request *service_user_register.EmailRequest) int {
	email := request.GetEmail()
	passwd := request.GetPassword()
	//检查邮箱和密码格式
	if !libraries.CheckEmail(email) {
		return codes.UserEmailInvalid
	}
	if !password.CheckPassword(passwd) {
		return codes.UserPasswordInvalid
	}
	return codes.Success
}
