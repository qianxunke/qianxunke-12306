package handler

import (
	book_handler "book-srv/handler/book-handler"
	"gitee.com/qianxunke/book-ticket-common/proto/task"
	"github.com/micro/go-micro/server"
	"log"
)

func Init() {
	book_handler.Init()
}

func RegisterHandler(server server.Server) {
	err := task.RegisterTaskServiceHandler(server, new(book_handler.Handler))
	if err != nil {
		log.Fatalf("[RegisterHandler] 注册 user_level handler 失败 ，%v", err)
		return
	}

}
