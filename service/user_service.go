package service

import (
	"GoChat/config"
	"GoChat/model"
	"log"
	"regexp"
)

var (
	db = config.PostgreSQL
)

//IsNameExist 判断用户名是否存在
func IsNameExist(name string) bool {
	user := model.User{}
	log.Printf("========用户名查询结果：%s", user.Username)
	result := db.Where("username = ?", name).Find(&user)
	log.Printf("用户名查询结果：%s,%d", user.Username, result.RowsAffected)
	return result.RowsAffected == 0
}

//CheckPhone 校验手机号
func CheckPhone(phone string) bool {
	phoneReg, _ := regexp.Match("^1(3\\d{2}|4[14-9]\\d|5([0-35689]\\d|7[1-79])|66\\d|7[2-35-8]\\d|8\\d{2}|9[13589]\\d)\\d{7}$", []byte(phone))
	if len(phone) != 11 { //判断手机号是否为11位
		return false
	} else if !phoneReg { //判断手机号是否符合正则
		return false
	}

	return true
}

//IsPhoneExist 判断电话号码是否存在
func IsPhoneExist(phone string) bool {
	user := model.User{}
	result := db.Where("phone = ?", phone).Find(&user)
	if result.RowsAffected != 0 {
		return false
	}
	return true
}
