package handler

import (
	"context"
	auth "gitee.com/qianxunke/book-ticket-common/proto/auth"
	"github.com/micro/go-micro/util/log"
	"surprise-shop-auth/model/access"
)

//声明service
var (
	accessService access.Service
)

func Init() {
	var err error
	accessService, err = access.GetService()
	if err != nil {
		log.Fatal("[Init] 初始化Handler错误，%s", err)
		return
	}
}

//实现proto接口
type Auth struct {
}

func (s *Auth) MakeAccessToken(ctx context.Context, req *auth.Request, rsp *auth.Response) (err error) {
	log.Log("[MakeAccessToken] 收到创建token请求: " + req.UserId + "  " + req.UserName)

	token, err := accessService.MakeAccessToken(&access.Subject{
		ID:   req.UserId,
		Name: req.UserName,
	})

	if err != nil {
		rsp.Error = &auth.Error{
			Message: err.Error(),
		}
		log.Logf("[MakeAccessToken] token生成失败，err：%s", err)
		return err
	}

	rsp.Token = token
	return
}

// DelUserAccessToken 清除用户token
func (s *Auth) DelUserAccessToken(ctx context.Context, req *auth.Request, rsp *auth.Response) error {
	log.Log("[DelUserAccessToken] 清除用户token")
	err := accessService.DelUserAccessToken(req.Token)
	if err != nil {
		rsp.Error = &auth.Error{
			Message: err.Error(),
		}

		log.Logf("[DelUserAccessToken] 清除用户token失败，err：%s", err)
		return err
	}

	return nil
}

// 鉴权用户
func (s *Auth) AuthenticationFromToken(ctx context.Context, req *auth.Request, rsp *auth.Response) error {
	userSub, err := accessService.AuthenticationFromToken(req.Token)
	if err != nil {
		rsp.Error = &auth.Error{
			Message: err.Error(),
		}
		rsp.Success = false
		log.Logf("[AuthenticationFromToken] 鉴权用户token失败，err：%s", err)
		return err
	}
	rsp.UserId = userSub.ID
	rsp.Success = true
	return nil
}
