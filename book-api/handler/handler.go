package handler

import (
	"book_ticket_api/handler/task_handler"
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/web"
	hystrixplugin "github.com/micro/go-plugins/wrapper/breaker/hystrix"
	"log"
)

func Init() {

}

//注册路由
func RegiserRouter(service web.Service, router *gin.Engine) {
	hystrix.DefaultTimeout = 5000
	sClient := hystrixplugin.NewClientWrapper()(service.Options().Service.Client())
	err := sClient.Init(
		client.Retries(0), //服务端错误请求重试1次
		client.Retry(func(ctx context.Context, req client.Request, retryCount int, err error) (bool, error) {
			//log.Log(req.Method(), retryCount, "请求重试： client retry")
			log.Println(req.Method(), retryCount, "请求重试： client retry")
			return true, nil
		}),
	)
	if err != nil {
		log.Fatal("[RegiserRouter] : %v", err)
	}
	inventoryRout := router.Group("/book")
	{
		productService := task_handler.Init(sClient)
		productRouter := inventoryRout.Group("/task")
		{
			productRouter.GET("/get/:taskId", productService.GetTask)
			productRouter.POST("/add", productService.AddTask)
			productRouter.POST("/list", productService.AddTask)
			productRouter.POST("/updateStatus", productService.UpdateTaskStatus)
			productRouter.GET("/userList", productService.GetUserTask)
		}

	}
}
