package contact

import (
	ws "GoChat/controller/websocket"
	"GoChat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Calculate(ctx *gin.Context) {
	for i := range ws.H.Clients {
		fmt.Println(i)
	}
	utils.Success(ctx, nil, strconv.Itoa(len(ws.H.Clients)))
}
