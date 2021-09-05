package main

import (
	"GoChat/config"
	"GoChat/router"
	"fmt"
)

var (
	WsClient interface{}
)

func main() {

	r := router.CollectRoute()
	port := fmt.Sprintf(":%d", config.Viper.ServerConfig.Port)
	panic(r.Run(port))
}
