package book

import (
	"book-srv/modules/book/book_dao"
	"book-srv/modules/book/book_service"
	"log"
)

func Init() {
	book_dao.Init()
	book_service.Init()
	s, err := book_service.GetService()
	if err != nil {
		log.Fatal(err.Error())
	}
	go s.StartBathTicket()
	go s.StartBathDoneError()
}
