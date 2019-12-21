package handler

import (
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/basic/api_common"
	"gitee.com/qianxunke/book-ticket-common/basic/common"
	"github.com/dgrijalva/jwt-go"
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
			claims, err := parseToken(ck)
			if err == nil {
				c.Request.Header.Add("userId", claims.Subject)
			}
		}
		c.Next()
	} else {
		ck := c.Request.Header.Get(common.RememberMeCookieName)
		if len(ck) == 0 {
			resonseEntity := &api_common.ResponseEntity{}
			resonseEntity.Message = "token为空，请先登陆!"
			resonseEntity.Code = http.StatusBadRequest
			c.JSON(http.StatusOK, resonseEntity)
			c.Abort()
			return
		}
		claims, err := parseToken(ck)
		//如果此token无效
		if err != nil {
			log.Logf("[AuthWrapper]，token不合法，无用户id")
			resonseEntity := &api_common.ResponseEntity{}
			resonseEntity.Message = "身份验证不通过，请先登陆!"
			resonseEntity.Code = http.StatusBadRequest
			c.JSON(http.StatusOK, resonseEntity)
			c.Abort()
			return
		}
		c.Request.Header.Add("userId", claims.Subject)
		c.Next()
	}
}

// parseToken 解析token
func parseToken(tk string) (c *jwt.StandardClaims, err error) {

	token, err := jwt.Parse(tk, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("不合法的token格式: %v", token.Header["alg"])
		}
		return []byte("I_AM_PASSWORD"), nil
	})

	// jwt 框架自带了一些检测，如过期，发布者错误等
	if err != nil {
		switch e := err.(type) {
		case *jwt.ValidationError:
			switch e.Errors {
			case jwt.ValidationErrorExpired:
				return nil, fmt.Errorf("[parseToken] 过期的token, err:%s", err)
			default:
				break
			}
			break
		default:
			break
		}

		return nil, fmt.Errorf("[parseToken] 不合法的token, err:%s", err)
	}

	// 检测合法
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("[parseToken] 不合法的token")
	}

	return mapClaimToJwClaim(claims), nil
}

// 把jwt的claim转成claims
func mapClaimToJwClaim(claims jwt.MapClaims) *jwt.StandardClaims {

	jC := &jwt.StandardClaims{
		Subject: claims["sub"].(string),
	}

	return jC
}
