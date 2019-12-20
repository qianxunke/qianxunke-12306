package sms

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	r "github.com/go-redis/redis"
	"net/http"
	"strconv"
	"time"
)

//短信通知类
var (
	smsapi = "http://api.smsbao.com/sms?"
	// 短信平台账号
	user = "qianxunke"
	// 短信平台密码
	password = "736567805"
	//短信签名
	sign = "【EsayGo】"
	// 要发送的短信内容
	content = "短信内容"
)

func main() {
	SendTicketSuccessInfoToUser("18334142052", "你好")
}

func SendTicketSuccessInfoToUser(phone string, userName string) {
	//https://api.smsbao.com/sms?u=USERNAME&p=PASSWORD&m=PHONE&c=CONTENT
	url := fmt.Sprintf("u=%s&p=%s&m=%s&c=%s",
		user,
		getMd5Pwd(password),
		phone,
		sign+"尊敬的"+userName+"：EsayGo已为您抢票成功，请在30分钟内到12306官方网站或APP完成支付。",
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

}

func SendRegisterMsm(code int64, phone string, rc *r.Client) (err error) {
	url := fmt.Sprintf("u=%s&p=%s&m=%s&c=%s",
		user,
		getMd5Pwd(password),
		phone,
		"您的验证码为"+strconv.FormatInt(code, 10)+"在3分钟内有效，请不要告诉任何人哦",
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
	redisErr := rc.Set(phone, code, time.Second*180).Err()
	if redisErr != nil {
		err = errors.New("Redis执行失败")
	}
	return
}

func getMd5Pwd(str string) (mdsString string) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	mdsString = hex.EncodeToString(md5Ctx.Sum(nil))
	return
}
