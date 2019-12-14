package subscriber

import (
	"book-user_srv/model/user_info"
	"context"
	"gitee.com/qianxunke/book-ticket-common/proto/user"
	"github.com/micro/go-micro/util/log"
)

var (
	uService user_info.UserService
)

func Init() {
	var err error
	uService, err = user_info.GetService()
	if err != nil {
		log.Fatal("[subscriber] productService : %v", err)
		return
	}
}

func UpdateUserPassengerHandle(ctx context.Context, msg *user.In_UpdateUserPassenger) error {
	log.Log("Handler Received message: ", msg)
	uService.UpdateUserPassenger(msg)
	return nil
}

func Handler(ctx context.Context, msg *user.Passenger) error {
	log.Log("Function Received message: ", "")
	return nil
}
