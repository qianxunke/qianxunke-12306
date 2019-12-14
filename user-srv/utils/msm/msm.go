package msm

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	r "github.com/go-redis/redis"
	"net/http"
	"strconv"
)

var (
	smsapi = "http://api.smsbao.com/sms?"
	// 短信平台账号
	user = "736765805"
	// 短信平台密码
	password = "736567805@qq.com"
	//短信签名
	sign = "【千寻客】"
	// 要发送的短信内容
	content = "短信内容"
)

func SendRegisterMsm(code int64, phone string, rc *r.Client) (err error) {
	url := fmt.Sprintf("u=%s&p=%s&m=%s&c=%s",
		user,
		GetMd5Pwd(password),
		phone,
		strconv.FormatInt(code, 10),
	)
	println("----" + smsapi + url)
	rsp, err := http.Get(smsapi + url)
	if err != nil {
		err = errors.New("短信发送失败")
		return
	}
	if rsp.StatusCode != http.StatusOK {
		err = errors.New("短信回调失败")
		return
	}
	//往redis写入验证码
	redisErr := rc.Do("SET", phone, code, "Ex", 30000).Err()
	if redisErr != nil {
		err = errors.New("Redis执行失败")
	}
	return
}

func GetMd5Pwd(str string) (mdsString string) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	mdsString = hex.EncodeToString(md5Ctx.Sum(nil))
	return
}
