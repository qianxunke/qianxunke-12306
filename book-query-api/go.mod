module book-query-api

go 1.13

replace gitee.com/qianxunke/book-ticket-common => ../book-ticket-common

require (
	gitee.com/qianxunke/book-ticket-common v0.0.0-00010101000000-000000000000
	gitee.com/qianxunke/surprise-shop-common v1.0.1 // indirect
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/gin-gonic/gin v1.5.0
	github.com/micro/go-micro v1.17.1
	github.com/micro/go-plugins v1.5.1
)
