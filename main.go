package main

import (
	"GoChat/config"
	"GoChat/router"
	"fmt"
)

func main() {
	r := router.CollectRoute()
	port := fmt.Sprintf(":%d", config.Viper.ServerConfig.Port)
	panic(r.Run(port))
}
