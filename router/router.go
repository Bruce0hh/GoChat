package router

import (
	"github.com/gin-gonic/gin"
)

// CollectRoute 添加router集合，前缀"/api"
func CollectRoute() *gin.Engine {

	r := gin.Default()
	apiRouter := r.Group("api")
	userRouter(apiRouter)

	return r

}

// UserRouter 添加用户相关api集合，前缀“/user”
func userRouter(r *gin.RouterGroup) {

	_ = r.Group("user")
	{

	}

}
