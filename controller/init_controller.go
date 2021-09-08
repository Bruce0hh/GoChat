package controller

import (
	"GoChat/config"
)

var (
	Db  = config.PostgreSQL
	Rdb = config.Redis
)
