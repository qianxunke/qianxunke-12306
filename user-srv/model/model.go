package model

import (
	"book-user_srv/global"
	"gitee.com/qianxunke/book-ticket-common/basic"
	"gitee.com/qianxunke/book-ticket-common/plugins/db"
	"gitee.com/qianxunke/book-ticket-common/plugins/redis"
	"gitee.com/qianxunke/book-ticket-common/proto/auth"
	"github.com/micro/go-micro/client"
)

func Init() {
	db.MasterEngine()
	global.RedisClient = redis.Redis()
	global.AuthClient = auth.NewAuthService(basic.AuthService, client.DefaultClient)

}
