package contact

import (
	"GoChat/utils"
	ws "GoChat/websocket"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func OnlineNumber() {

}

func Calculate(ctx *gin.Context) {
	for i := range ws.H.Clients {
		fmt.Println(i)
	}
	utils.Success(ctx, nil, strconv.Itoa(len(ws.H.Clients)))
}
