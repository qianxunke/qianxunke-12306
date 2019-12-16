package main

import (
	"book-user_api/handler"
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/basic"
	"gitee.com/qianxunke/book-ticket-common/basic/common"
	"gitee.com/qianxunke/book-ticket-common/basic/config"
	"gitee.com/qianxunke/book-ticket-common/basic/lib/tracer"
	"gitee.com/qianxunke/book-ticket-common/basic/lib/wrapper/tracer/opentracing/gin2micro"
	"github.com/gin-gonic/gin"
	"github.com/micro/cli"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/config/source/grpc"
	"github.com/micro/go-plugins/registry/consul"
	"github.com/opentracing/opentracing-go"
	"time"
)

var (
	appName = "user_api"
	cfg     = &appCfg{}
)

type appCfg struct {
	common.AppCfg
}

func main() {
	initCfg()
	micReg := consul.NewRegistry(registryOptions)
	gin2micro.SetSamplingFrequency(50)
	t, io, err := tracer.NewTracer(cfg.Name, "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	service := web.NewService(
		web.Name(cfg.Name),
		web.Registry(micReg),
		web.Version(cfg.Version),
		web.RegisterTTL(time.Second*15), //健康检查，15秒后重新向注册中心注册
		web.RegisterInterval(time.Second*10),
		web.Address(cfg.Addr()),
	)
	// initialise service
	if err := service.Init(web.Action(func(i *cli.Context) {
		handler.Init()
	})); err != nil {
		log.Fatal(err)
	}

	router := gin.New()
	router.Use(gin2micro.TracerWrapper)
	router.Use(handler.AuthWrapper)
	handler.RegiserRouter(service, router)

	service.Handle("/", router)
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func registryOptions(ops *registry.Options) {
	consulCfg := &common.Consul{}
	err := config.C().App("consul", consulCfg)
	if err != nil {
		panic(err)
	}
	ops.Addrs = []string{fmt.Sprintf("%s:%d", consulCfg.Host, consulCfg.Port)}
}

func initCfg() {
	if common.WhileList == nil {
		common.WhileList = &common.UrlWhileList{}
	}

	source := grpc.NewSource(
		grpc.WithAddress(common.ControlCenterAddress),
		grpc.WithPath(appName),
	)

	basic.Init(config.WithSource(source))
	err := config.C().App(appName+"_app", cfg)
	if err != nil {
		panic(err)
	}
	err = config.C().App("url_whitelist", &common.WhileList.Url_whitelist)
	if err != nil {
		panic(err)
	}
	log.Logf("[initCfg] 配置，cfg：%v", common.WhileList)
	return
}
