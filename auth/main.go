package main

import (
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/basic"
	"gitee.com/qianxunke/book-ticket-common/basic/common"
	"gitee.com/qianxunke/book-ticket-common/basic/config"
	"gitee.com/qianxunke/book-ticket-common/basic/lib/tracer"
	auth "gitee.com/qianxunke/book-ticket-common/proto/auth"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/config/source/grpc"
	"github.com/micro/go-plugins/registry/consul"
	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"surprise-shop-auth/handler"
	"surprise-shop-auth/model"
	"time"
)

var (
	appName = "auth_srv"
	cfg     = &appCfg{}
)

type appCfg struct {
	common.AppCfg
}

func main() {
	initCfg()
	// 使用consul注册
	micReg := consul.NewRegistry(registryOptions)
	t, io, err := tracer.NewTracer(cfg.Name, "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)
	// New Service
	service := micro.NewService(
		micro.Name(cfg.Name),
		micro.Registry(micReg),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*15), //健康检查，15秒后重新向注册中心注册
		micro.RegisterInterval(time.Second*10),
		micro.WrapHandler(ocplugin.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	// Initialise service
	service.Init(
		micro.Action(func(context *cli.Context) {
			model.Init()
			handler.Init()
		}),
	)
	// Register Handler
	_ = auth.RegisterAuthHandler(service.Server(), new(handler.Auth))
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
	log.Logf("[initCfg] 配置，cfg：%v", cfg)
	return
}
