package jwt

import "gitee.com/qianxunke/surprise-shop-common/basic"

// Jwt 配置 接口
type Jwt struct {
	SecretKey string `json:"secretKey"`
}

// init 初始化Redis
func init() {
	basic.Register(initJwt)
}

func initJwt() {

}
