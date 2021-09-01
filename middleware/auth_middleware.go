package middleware

import (
	"GoChat/config"
	"GoChat/model"
	"GoChat/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var (
	db   = config.PostgreSQL
	user model.User
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		//校验token
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			utils.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足！")
			//中间件中断需要Abort
			ctx.Abort()
			return
		}
		tokenString = tokenString[7:]

		token, claims, err := config.ParseToken(tokenString)
		if err != nil || !token.Valid {
			utils.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足！")
			ctx.Abort()
			return
		}

		//验证通过则获取Claims的UserId
		userId := claims.UserId
		result := db.First(&user, userId)

		//若数据库不存在用户信息
		if result.RowsAffected == 0 {
			utils.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足！")
			ctx.Abort()
			return
		}
		//用户存在
		ctx.Set("user", user)
		//跳转下个中间件
		ctx.Next()
	}
}
