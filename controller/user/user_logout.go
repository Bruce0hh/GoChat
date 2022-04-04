package user

import (
	"GoChat/controller"
	"GoChat/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
)

func Logout(ctx *gin.Context) {

	username, _ := ctx.GetQuery("username")
	val, err := controller.Rdb.Get(ctx, username).Result()

	switch {
	case err == redis.Nil:
		utils.Response(ctx, http.StatusInternalServerError, 500, nil, "key does not exist")
		return
	case err != nil:
		utils.Response(ctx, http.StatusInternalServerError, 500, nil, "Get failed")
		return
	case val == "":
		utils.Response(ctx, http.StatusInternalServerError, 500, nil, "value is empty")
		return
	}
	log.Println("========>" + username)
	controller.Rdb.Del(ctx, username)
	utils.Success(ctx, gin.H{"user": username}, "该账号已下线")

}
