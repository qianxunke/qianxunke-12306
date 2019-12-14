package global

import (
	auth "gitee.com/qianxunke/book-ticket-common/proto/auth"
	r "github.com/go-redis/redis"
)

var (
	RedisClient *r.Client
	AuthClient  auth.AuthService
)
