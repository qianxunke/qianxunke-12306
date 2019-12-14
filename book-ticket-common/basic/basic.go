package basic

import (
	"gitee.com/qianxunke/surprise-shop-common/basic/config"
)

var (
	pluginFuncs []func()
)

type Options struct {
	EnableDB    bool
	EnableRedis bool
	cfgOps      []config.Option
}

type Option func(o *Options)

func Init(opts ...config.Option) {
	// 初始化配置
	config.Init(opts...)

	// 加载依赖配置的插件
	for _, f := range pluginFuncs {
		f()
	}
}

func Register(f func()) {
	pluginFuncs = append(pluginFuncs, f)
}
