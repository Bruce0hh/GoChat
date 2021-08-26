package config

import (
	"GoChat/config/model"
	"github.com/spf13/viper"
)

// InitViperConfig 初始化viper
func InitViperConfig() model.GlobalConfig {

	//实例化viper
	v := viper.New()
	//viper读取配置文件
	v.SetConfigFile("config/application.yml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//给serverConfig赋值
	globalConfig := model.GlobalConfig{}
	if err := v.Unmarshal(&globalConfig); err != nil {
		panic(err)
	}
	return globalConfig

}
