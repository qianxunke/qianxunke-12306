package http_util

import (
	"encoding/json"
	"net/http"
	"qianxunke-12306/modules/conversation"
	"strings"
)

/**
 * 伪造请求头的方法
 */
func SetReqHeader(httpReq *http.Request) {
	httpReq.Header.Set("User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.26 Safari/537.36 Core/1.63.6788.400 QQBrowser/10.3.2864.400")
	httpReq.Header.Set("Host", "kyfw.12306.cn")
	httpReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	httpReq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
}

//向请求添加cookie
func AddReqCookie(cs []*http.Cookie, httpReq *http.Request) {
	for _, c := range cs {
		httpReq.AddCookie(c)
	}
}

func CookieChange(conversation *conversation.Conversation, cs []*http.Cookie) {
	if len(cs) <= 0 {
		return
	}
	if len(conversation.C) <= 0 {
		conversation.C = cs
		return
	}
	for _, c := range cs {
		have := false
		for i := 0; i < len(conversation.C); i++ {
			if c.Name == conversation.C[i].Name {
				conversation.C[i] = c
				have = true
				break
			}
		}
		if !have {
			conversation.C = append(conversation.C, c)
		}
	}
}

func TimeStrapStringJsonToBean(str string, bean interface{}) (err error) {
	rs := []rune(str)
	start := strings.Index(str, "(")
	nRs := []byte(string(rs[start+1 : len(rs)-2]))
	err = json.Unmarshal(nRs, &bean)
	return
}
func StringJsonToBean(str string, bean *interface{}) (err error) {
	nRs := []byte(str)
	err = json.Unmarshal(nRs, &bean)
	return
}
