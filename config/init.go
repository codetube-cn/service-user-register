package config

// ServiceConfig 网关配置
var ServiceConfig *Config

func init() {
	//初始化配置
	ServiceConfig = InitConfig()
}
