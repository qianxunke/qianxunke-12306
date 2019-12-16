package m_client

import (
	"gitee.com/qianxunke/book-ticket-common/basic"
	"gitee.com/qianxunke/book-ticket-common/proto/auth"
	"github.com/micro/go-micro/client"
)

var (
	AuthClient auth.AuthService
)

func Init() {
	AuthClient = auth.NewAuthService(basic.AuthService, client.DefaultClient)
}
