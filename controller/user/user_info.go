package user

import (
	"GoChat/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Info todo：现在是无需传参，根据 Bearer Token 获取信息，后续需要控制每个用户的查询操作
func Info(ctx *gin.Context) {
	if user, ok := ctx.Get("user"); ok {
		utils.Success(ctx, gin.H{"user": user}, "成功返回用户信息！")
		return
	}
	utils.Response(ctx, http.StatusInternalServerError, 500, nil, "查无此人")

}
