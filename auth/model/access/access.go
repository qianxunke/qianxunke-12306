package access

import (
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/basic/config"
	"gitee.com/qianxunke/book-ticket-common/plugins/jwt"
	"github.com/micro/go-micro/util/log"
	"sync"

	"gitee.com/qianxunke/book-ticket-common/plugins/redis"
	r "github.com/go-redis/redis"
)

//到时候会有人实现你的哈
type service struct {
}

var (
	s   *service
	ca  *r.Client
	m   sync.RWMutex
	cfg = &jwt.Jwt{}
)

//接口
type Service interface {
	//生成toke，
	MakeAccessToken(subject *Subject) (ret string, err error)

	//得到缓存的token
	GetCacheAccessToken(subject *Subject) (ret string, err error)

	//清除用户token
	DelUserAccessToken(token string) (err error)

	//解析token获取用户信息
	AuthenticationFromToken(tk string) (subject *Subject, err error)
}

//获取服务
func GetService() (Service, error) {
	if s == nil {
		return nil, fmt.Errorf("[GetService] GetService 未初始化")
	}
	return s, nil
}

func Init() {
	m.Lock()
	defer m.Unlock()
	if s != nil {
		return
	}
	err := config.C().App("jwt", cfg)
	if err != nil {
		panic(err)
	}

	log.Logf("[initCfg] 配置，cfg：%v", cfg)

	ca = redis.Redis()

	s = &service{}
}
