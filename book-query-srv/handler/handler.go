package handler

import (
	"book-query-srv/handler/train"
	"gitee.com/qianxunke/book-ticket-common/proto/ticket"
	"github.com/micro/go-micro/server"
	"log"
)

func Init() {
	train.Init()
}

func RegisterHandler(server server.Server) {
	err := ticket.RegisterTrainServiceHandler(server, new(train.Handler))
	if err != nil {
		log.Fatalf("[RegisterHandler] 注册 ticket handler 失败 ，%v", err)
		return
	}

}
