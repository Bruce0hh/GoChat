package config

import (
	"GoChat/config/model"
	"fmt"
	"github.com/spf13/viper"
)

var (
	Viper = InitViper()
)

// InitViper 初始化viper
func InitViper() model.GlobalConfig {

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
	fmt.Println(globalConfig.ServerConfig.Port)
	return globalConfig

}
