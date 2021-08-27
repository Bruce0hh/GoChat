package controller

import (
	"GoChat/config"
	"GoChat/model"
)

var (
	user model.User
	db   = config.PostgreSQL
	rdb  = config.Redis
)
