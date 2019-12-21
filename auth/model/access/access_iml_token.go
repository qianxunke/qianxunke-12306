package access

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/util/log"
	"time"
)

var (
	// tokenExpiredDate app token过期日期 7天
	tokenExpiredDate = 3600 * 24 * 7 * time.Second
	// tokenIDKeyPrefix tokenID 前缀
	tokenIDKeyPrefix  = "token:auth:id:"
	tokenExpiredTopic = "com.surprise.shop.topic.auth.tokenExpired"
)

//token 持有者
type Subject struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// standardClaims token 标准的Claims
type standardClaims struct {
	SubjectID string `json:"subjectId,omitempty"`
	Name      string `json:"name,omitempty"`
	jwt.StandardClaims
}

//生成token并保存到redis
func (s *service) MakeAccessToken(subject *Subject) (ret string, err error) {

	m, err := s.createTokenClaims(subject)
	if err != nil {
		return "", fmt.Errorf("[MakeAccessToken] 创建token Claim 失败，err: %s", err)
	}
	//创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, m)
	ret, err = token.SignedString([]byte(cfg.SecretKey))

	if err != nil {
		return "", fmt.Errorf("[MakeAccessToken] 创建token失败，err: %s", err)
	}
	//保存加密的token到redis
	err = s.saveTokenToCache(subject, ret)
	if err != nil {
		return "", fmt.Errorf("[MakeAccessToken] 保存token到缓存失败，err: %s", err)
	}

	return

}

// GetCachedAccessToken 获取token
func (s *service) GetCacheAccessToken(subject *Subject) (ret string, err error) {
	ret, err = s.getTokenFromCache(subject)
	if err != nil {
		return "", fmt.Errorf("[GetCachedAccessToken] 从缓存获取token失败，err: %s", err)
	}

	return
}

//清除用户toekn
func (s *service) DelUserAccessToken(token string) (err error) {
	//解析token
	claims, err := s.parseToken(token)
	if err != nil {
		return fmt.Errorf("[DelUserAccessToken] 错误的token，err: %s", err)
	}
	//通过解析到的用户id删除
	err = s.delTokenFromCache(&Subject{
		ID: claims.Subject,
	})
	if err != nil {
		return fmt.Errorf("[DelUserAccessToken] 清除用户token，err: %s", err)
	}
	//广播删除
	msg := &broker.Message{
		Body: []byte(claims.Subject),
	}
	if err := broker.Publish(tokenExpiredTopic, msg); err != nil {
		log.Logf("[pub] 发布消息失败： %v", err)
	} else {
		fmt.Println("[pub] 发布消息：", string(msg.Body))
	}
	return
}

func (s *service) AuthenticationFromToken(tk string) (subject *Subject, err error) {
	claims, err := s.parseToken(tk)
	//如果此token无效
	if err != nil {
		return
	}
	subject = &Subject{
		ID: claims.Subject,
	}
	/*
		cacheToken, err := s.getTokenFromCache(subject)
		if err != nil || len(cacheToken) == 0 || cacheToken != tk {
			return nil, fmt.Errorf("[AuthenticationFromToken] 从缓存获取token失败，err: %s", err)
		}

	*/
	return
}
