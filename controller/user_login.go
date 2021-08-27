package controller

import (
	"GoChat/config"
	"GoChat/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func Login(ctx *gin.Context) {

	//获取参数
	name := ctx.PostForm("username")
	password := ctx.PostForm("password")

	//判断用户是否存在
	if service.IsNameExist(name) {
		config.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在！")
		return
	}
	db.Where("username = ?", name).Find(&user)
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		config.Response(ctx, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}

	//发放Token
	token, err := config.ReleaseToken(user)
	if err != nil {
		config.Response(ctx, http.StatusInternalServerError, 500, nil, "生成token失败")
		log.Printf("token generate error : %v\n", err)
		return
	}

	//登录成功
	config.Success(ctx, gin.H{"token": token}, "登录成功")
	return
}
