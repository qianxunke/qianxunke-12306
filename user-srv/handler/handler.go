package handler

import (
	"book-user_srv/handler/user_info"
	"gitee.com/qianxunke/book-ticket-common/proto/user"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/util/log"
)

func Init() {
	user_info.Init()

}
func RegisterHandler(server server.Server) {
	err := user.RegisterUserInfoHandler(server, new(user_info.Handler))
	if err != nil {
		log.Fatalf("[RegisterHandler] 注册 user_level handler 失败 ，%v", err)
		return
	}

}
