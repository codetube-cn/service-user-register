package server

import (
	"codetube.cn/core/codes"
	core_libraries "codetube.cn/core/libraries"
	service_user_register "codetube.cn/proto/service-user-register"
	"codetube.cn/service-user-register/components"
	"codetube.cn/service-user-register/libraries"
	"codetube.cn/service-user-register/libraries/check"
	"codetube.cn/service-user-register/models"
	"context"
	"log"
	"strconv"
)

// Mobile 使用邮箱注册
func (s *UserRegisterServer) Mobile(c context.Context, request *service_user_register.MobileRequest) (*service_user_register.RegisterResultResponse, error) {
	message := "success"
	status := s.checkMobileParams(request)
	if status != codes.Success {
		return &service_user_register.RegisterResultResponse{
			Status:  int64(status),
			Message: message, //@todo 数字转文字
			Id:      0,
		}, nil
	}

	mobile := request.GetMobile()
	id := int64(0)
	//检查手机号是否重复
	userExist := check.UserExistByMobile(mobile)
	if userExist {
		status = codes.UserMobileIsExist
	} else {
		//写入数据
		//定义一个用户，并初始化数据
		user := models.User{
			Mobile:  mobile,
			Enabled: 1,
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
		//清除验证码
		cacheKey := libraries.GetMobileRegisterVerifyCodeCacheKey(mobile)
		_, err := components.CommonRedis.Del(context.Background(), cacheKey).Result()
		if err != nil {
			log.Println("用户手机号注册从缓存删除验证码失败：", cacheKey, err)
		}
	}
	return &service_user_register.RegisterResultResponse{
		Status:  int64(status),
		Message: message,
		Id:      id,
	}, nil
}

func (s *UserRegisterServer) checkMobileParams(request *service_user_register.MobileRequest) int {
	mobile := request.GetMobile()
	verifyCode := request.GetVerifyCode()
	//检查手机号格式
	if !core_libraries.CheckMobile(mobile) {
		return codes.UserMobileInvalid
	}
	if len(verifyCode) < 4 || len(verifyCode) > 8 {
		return codes.VerifyCodeInvalid
	}
	//检查验证码是否正确
	cacheKey := libraries.GetMobileRegisterVerifyCodeCacheKey(mobile)
	cacheVerifyCode, err := components.CommonRedis.Get(context.Background(), cacheKey).Result()
	if err != nil {
		log.Println("用户手机号注册从缓存获取验证码失败：", cacheKey, err)
		return codes.VerifyCodeInvalid
	}
	if cacheVerifyCode != verifyCode {
		return codes.VerifyCodeInvalid
	}
	return codes.Success
}
