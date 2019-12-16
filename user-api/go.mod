module book-user_api

go 1.13

require (
	gitee.com/qianxunke/book-ticket-common v0.0.0-00010101000000-000000000000
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/gin-gonic/gin v1.5.0
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.17.1
	github.com/micro/go-plugins v1.5.1
	github.com/opentracing/opentracing-go v1.1.0
)

replace gitee.com/qianxunke/book-ticket-common => ../book-ticket-common
