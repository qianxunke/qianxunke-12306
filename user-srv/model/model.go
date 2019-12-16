package model

import (
	"book-user_srv/global"
	"book-user_srv/model/user_info"
	"gitee.com/qianxunke/book-ticket-common/basic"
	"gitee.com/qianxunke/book-ticket-common/plugins/redis"
	"gitee.com/qianxunke/book-ticket-common/proto/auth"
	"github.com/micro/go-micro/client"
)

func Init() {
	user_info.Init()
	global.RedisClient = redis.Redis()
	global.AuthClient = auth.NewAuthService(basic.AuthService, client.DefaultClient)

}
