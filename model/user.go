package model

import (
	"gorm.io/gorm"
)

//User 用户表
type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`

	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Status    string `json:"status"`
	Sex       string `json:"sex"`
	Signature string `json:"signature"`
}
