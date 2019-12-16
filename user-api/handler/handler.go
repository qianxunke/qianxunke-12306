package handler

import (
	"book-user_api/handler/user_info"
	"book-user_api/m_client"
	"context"
	"gitee.com/qianxunke/book-ticket-common/basic"
	auth "gitee.com/qianxunke/book-ticket-common/proto/auth"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/web"
	hystrixplugin "github.com/micro/go-plugins/wrapper/breaker/hystrix"
)

func Init() {
	m_client.AuthClient = auth.NewAuthService(basic.AuthService, client.DefaultClient)
}

//注册路由
func RegiserRouter(service web.Service, router *gin.Engine) {
	hystrix.DefaultTimeout = 5000
	sClient := hystrixplugin.NewClientWrapper()(service.Options().Service.Client())
	err := sClient.Init(
		client.Retries(1), //服务端错误请求重试1次
		client.Retry(func(ctx context.Context, req client.Request, retryCount int, err error) (bool, error) {
			log.Log(req.Method(), retryCount, "请求重试： client retry")
			return true, nil
		}),
	)
	if err != nil {
		log.Fatal("[RegiserRouter] : %v", err)
	}
	userRout := router.Group("/user")
	{
		apiService := user_info.Init(sClient)
		userRout.POST("/login", apiService.Login)
		userRout.POST("/out", apiService.Logout)
		userRout.POST("/register", apiService.Register)
		userRout.POST("/code", apiService.GetCode)
		userRout.POST("/list", apiService.GetUserInfoList)
		userRout.GET("/info", apiService.GetUserInfo)
		userRout.POST("/login12306", apiService.Login12306)

	}
	//user_level

}
