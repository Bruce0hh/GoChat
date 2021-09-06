package config

import (
	"github.com/spf13/viper"
)

var (
	Viper      = initViper()
	PostgreSQL = initPostgres()
	Redis      = initRedis()
	MongoDB    = initMongoDB()
)

// InitViper 初始化viper
func initViper() GlobalConfig {

	//实例化viper
	v := viper.New()
	//viper读取配置文件
	v.SetConfigFile("config/application.yml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//给serverConfig赋值
	globalConfig := GlobalConfig{}
	if err := v.Unmarshal(&globalConfig); err != nil {
		panic(err)
	}

	return globalConfig

}
