package libraries

//GetMobileRegisterVerifyCodeCacheKey 获取手机号注册时验证码缓存 key
func GetMobileRegisterVerifyCodeCacheKey(mobile string) string {
	return "user_mobile_register_verify_code_" + mobile
}
