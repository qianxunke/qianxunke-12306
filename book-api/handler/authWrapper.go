package handler

import (
	"book_ticket_api/m_client"
	"context"
	"gitee.com/qianxunke/book-ticket-common/basic/api_common"
	"gitee.com/qianxunke/book-ticket-common/basic/common"
	"gitee.com/qianxunke/book-ticket-common/proto/auth"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/util/log"
	"net/http"
)

//token 持有者
type Subject struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// AuthWrapper 认证wrapper
func AuthWrapper(c *gin.Context) {
	log.Logf("[AuthWrapper]:请求 URL: %v", c.Request.RequestURI)
	if common.WhileList.IsInWileList(c.Request.RequestURI) {
		//假如该请求是登陆后
		ck := c.Request.Header.Get(common.RememberMeCookieName)
		// token不存在，则状态异常，无权限
		if len(ck) > 0 {
			tokenRsp, err := m_client.AuthClient.AuthenticationFromToken(context.TODO(), &auth.Request{
				Token: ck,
			})
			if err == nil && tokenRsp.Success {
				c.Request.Header.Add("userId", tokenRsp.UserId)
			}
		}
		c.Next()
	} else {
		ck := c.Request.Header.Get(common.RememberMeCookieName)
		if len(ck) == 0 {
			resonseEntity := &api_common.ResponseEntity{}
			resonseEntity.Message = "身份验证不通过，请先登陆!"
			resonseEntity.Code = http.StatusBadRequest
			c.JSON(http.StatusBadRequest, resonseEntity)
			c.Abort()
			return
		}
		tokenRsp, err := m_client.AuthClient.AuthenticationFromToken(context.TODO(), &auth.Request{
			Token: ck,
		})
		if err == nil && tokenRsp.Success {
			c.Request.Header.Add("userId", tokenRsp.UserId)
			c.Next()
		} else {
			log.Logf("[AuthWrapper]，token不合法，无用户id")
			resonseEntity := &api_common.ResponseEntity{}
			resonseEntity.Message = "身份验证不通过，请先登陆!"
			resonseEntity.Code = http.StatusBadRequest
			c.JSON(http.StatusBadRequest, resonseEntity)
			c.Abort()
			return
		}
	}
}
