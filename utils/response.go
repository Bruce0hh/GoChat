package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(ctx *gin.Context, httpStatus int, code int, data gin.H, msg string) {
	ctx.JSON(httpStatus, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

func Success(ctx *gin.Context, data gin.H, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
		"msg":  msg,
	})
}

func Fail(ctx *gin.Context, code int, data gin.H, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": data,
		"msg":  msg,
	})
}
