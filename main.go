package main

import (
	"GoChat/config"
	"GoChat/router"
	ws "GoChat/websocket"
	"fmt"
)

var (
	WsClient interface{}
)

func main() {

	WsClient = ws.CreateHubFactory()
	if wh, ok := WsClient.(*ws.Hub); ok {
		go wh.Run()
	}
	r := router.CollectRoute()
	port := fmt.Sprintf(":%d", config.Viper.ServerConfig.Port)
	panic(r.Run(port))
}
