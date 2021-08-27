package controller

import (
	"GoChat/config"
	"GoChat/model"
	"GoChat/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

//Register 用户名+密码
func Register(ctx *gin.Context) {

	//获取参数
	name := ctx.PostForm("username")
	password := ctx.PostForm("password")

	//用户名校验
	if len(name) == 0 {
		config.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "姓名不能为空！")
		return
	}

	if !service.IsNameExist(name) {
		config.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "该用户名已被注册！")
		return
	}

	//todo:密码正则校验特殊字符等
	//密码校验
	if len(password) < 6 || len(password) > 12 {
		config.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码长度在6~12位之间！")
		return
	}

	//Bcrypt加密密码
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		config.Response(ctx, http.StatusInternalServerError, 500, nil, "Bcrypt加密失败！")
		return
	}

	//创建用户
	newUser := model.User{
		Model:    gorm.Model{},
		Username: name,
		Password: string(bcryptPassword),
	}
	db.Create(&newUser)

	//返回结果
	config.Success(ctx, nil, "注册成功")
	return
}
