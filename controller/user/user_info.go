package user

import (
	"GoChat/utils"
	"github.com/gin-gonic/gin"
)

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	utils.Success(ctx, gin.H{"user": user}, "成功返回用户信息！")
}
