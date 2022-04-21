package check

import (
	"codetube.cn/service-user-register/components"
	"codetube.cn/service-user-register/models"
)

// UserExistByUsername 检查用户名是否存在
func UserExistByUsername(username string) bool {
	user := &models.User{}
	components.UserDB.Where("username = ?", username).First(user)
	if user.ID < 1 {
		return false
	}
	return true
}

// UserExistByEmail 检查邮箱是否存在
func UserExistByEmail(username string) bool {
	user := &models.User{}
	components.UserDB.Where("email = ?", username).First(user)
	if user.ID < 1 {
		return false
	}
	return true
}
