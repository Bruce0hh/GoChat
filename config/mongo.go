package config

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func initMongoDB() *mongo.Database {
	host := Viper.MongoDBConfig.Host
	port := Viper.MongoDBConfig.Port
	database := Viper.MongoDBConfig.Database

	uri := "mongodb://%s:%d/%s"
	uri = fmt.Sprintf(uri, host, port, database)
	// 连接客户端配置
	clientOptions := options.Client().ApplyURI(uri)
	// 连接到 MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		zap.Error(err)
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		zap.Error(err)
	}

	return client.Database("go-chat")
}
