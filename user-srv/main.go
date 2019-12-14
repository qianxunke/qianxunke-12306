package main

import (
	"book-user_srv/handler"
	"book-user_srv/model"
	"book-user_srv/subscriber"
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/basic"
	"gitee.com/qianxunke/book-ticket-common/basic/common"
	"gitee.com/qianxunke/book-ticket-common/basic/config"
	"gitee.com/qianxunke/book-ticket-common/basic/lib/tracer"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/config/source/grpc"
	"github.com/micro/go-plugins/registry/consul"
	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"
	opentracing "github.com/opentracing/opentracing-go"

	"time"
	//"surprise/user-srv/subscriber"
)

var (
	appName = "user_srv"
	cfg     = &appCfg{}
)

type appCfg struct {
	common.AppCfg
}

func main() {
	//初始化
	initCfg()
	microReg := consul.NewRegistry(registryOptions)
	t, io, err := tracer.NewTracer(cfg.Name, "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// New Service
	service := micro.NewService(
		micro.Name(cfg.Name),
		micro.Registry(microReg),
		micro.Version(cfg.Version),
		micro.RegisterTTL(time.Second*15),
		micro.RegisterInterval(time.Second*10),
		micro.WrapHandler(ocplugin.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	service.Init(
		micro.Action(func(context *cli.Context) {
			model.Init()
			handler.Init()
			subscriber.Init()
		}))
	//注册handler
	handler.RegisterHandler(service.Server())

	err = micro.RegisterSubscriber(common.Topic_Product_Es_Add, service.Server(), subscriber.UpdateUserPassengerHandle)
	if err != nil {
		log.Fatal(err)
	}
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

	err := config.C().App(appName+"_app", cfg) //获取
	if err != nil {
		panic(err)
	}
	log.Logf(appName+" [initCfg] 配置，cfg：%v", cfg)
	return
}
