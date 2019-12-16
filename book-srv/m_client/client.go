package m_client

import (
	"gitee.com/qianxunke/book-ticket-common/basic"
	user "gitee.com/qianxunke/book-ticket-common/proto/user"
	"github.com/micro/go-micro/client"
)

func Init() {
	UserClient = user.NewUserInfoService(basic.UserService, client.DefaultClient)
}

var (
	UserClient user.UserInfoService
)
