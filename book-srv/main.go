package main

import (
	"book-srv/handler"
	"book-srv/m_client"
	"book-srv/modules"
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/basic"
	"gitee.com/qianxunke/book-ticket-common/basic/common"
	"gitee.com/qianxunke/book-ticket-common/basic/config"
	"gitee.com/qianxunke/book-ticket-common/basic/lib/tracer"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/config/source/grpc"
	"github.com/micro/go-plugins/registry/consul"
	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"log"
	"time"
)

var (
	appName = "book_srv"
	cfg     = &userCfg{}
)

type userCfg struct {
	common.AppCfg
}

////go run main.go plugin.go --broker=nats --broker_address=127.0.0.1:4222
func main() {
	initCfg()
	t, io, err := tracer.NewTracer(cfg.Name, "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	service := micro.NewService(
		micro.Name(cfg.Name),
		micro.Registry(consul.NewRegistry(registryOptions)),
		micro.Version(cfg.Version),
		micro.RegisterTTL(time.Second*15), //健康检查，15秒后重新向注册中心注册
		micro.RegisterInterval(time.Second*10),
		micro.WrapHandler(ocplugin.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	// Initialise service
	service.Init(micro.Action(func(context *cli.Context) {
		m_client.Init()
		modules.Init()
		handler.Init()

	}))
	//注册handler
	handler.RegisterHandler(service.Server())

	// Run service
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
	source := grpc.NewSource(
		grpc.WithAddress(common.ControlCenterAddress),
		grpc.WithPath(appName),
	)

	basic.Init(config.WithSource(source))

	err := config.C().App(appName+"_app", cfg)
	if err != nil {
		panic(err)
	}
	log.Printf("[initCfg] 配置，cfg：%v", cfg)
	return
}
