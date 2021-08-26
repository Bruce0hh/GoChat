package config

import (
	"GoChat/model"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initPostgres() *gorm.DB {

	//从 viper 获取 postgres 配置信息
	postgreSQLConfig := Viper.PostgreSQLConfig
	host := postgreSQLConfig.Host
	port := postgreSQLConfig.Port
	username := postgreSQLConfig.Name
	password := postgreSQLConfig.Password
	database := postgreSQLConfig.Database

	//拼接 postgres 的数据库连接
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		host, port, username, password, database)

	//关闭 prepared statement 缓存
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: false,
	}), &gorm.Config{})

	if err != nil {
		panic("failed to connect database, err:" + err.Error())
	}
	err1 := db.AutoMigrate(&model.User{})
	if err1 != nil {
		return nil
	}
	return db
}
