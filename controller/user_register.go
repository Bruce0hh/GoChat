package controller

import (
	"GoChat/config"
	"GoChat/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"regexp"
)

func Register(ctx *gin.Context) {

	//获取参数
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")
	phone := ctx.PostForm("phone")

	//数据校验
	if len(name) == 0 {
		config.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "姓名不能为空！")
		return
	}

	//密码校验
	if len(password) < 6 || len(password) > 12 {
		config.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码长度在6~12位之间！")
		return
	}

	//手机号校验
	phoneReg, _ := regexp.Match("^1(3\\d{2}|4[14-9]\\d|5([0-35689]\\d|7[1-79])|66\\d|7[2-35-8]\\d|8\\d{2}|9[13589]\\d)\\d{7}$", []byte(phone))
	if len(phone) != 11 {
		config.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机长度必须为11！")
		return
	} else if !phoneReg {
		config.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "非法手机号！")
		return
	}
	if isPhoneExist(db, phone) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "用户已存在",
		})
		config.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已存在！")
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
		Phone:    phone,
	}
	db.Create(&newUser)

	//返回结果
	config.Success(ctx, nil, "注册成功")
	return
}

//判断电话号码是否存在
func isPhoneExist(db *gorm.DB, phone string) bool {

	result := db.Where("phone = ?", phone).Find(&user)
	if result.RowsAffected != 0 {
		return true
	}
	return false
}
