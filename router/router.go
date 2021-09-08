package router

import (
	"GoChat/controller/contact"
	"GoChat/controller/user"
	"GoChat/middleware"
	"GoChat/websocket"
	"github.com/gin-gonic/gin"
)

// CollectRoute 添加router集合，前缀"/api"
func CollectRoute() *gin.Engine {

	r := gin.Default()
	apiRouter := r.Group("api")
	wsRoute := r.Group("ws")
	userRouter(apiRouter)
	chatRouter(wsRoute)

	return r

}

// userRouter 添加用户相关api集合，前缀“/user”
func userRouter(r *gin.RouterGroup) {
	userRouter := r.Group("user")
	{
		userRouter.POST("register", user.Register)
		userRouter.POST("login", user.Login)
		userRouter.GET("info", middleware.AuthMiddleWare(), user.Info)
	}

}

func chatRouter(r *gin.RouterGroup) {
	chatRouter := r.Group("chat")
	{
		chatRouter.GET("", websocket.Chat)
		chatRouter.GET("calculate", contact.Calculate)
	}
}
