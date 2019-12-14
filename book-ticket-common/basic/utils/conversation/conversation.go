package conversation

import "net/http"

//会话
type Conversation struct {
	C      []*http.Cookie //请求cokkie
	Client *http.Client   //客户端
	UserId int64          //哪个账户
}
