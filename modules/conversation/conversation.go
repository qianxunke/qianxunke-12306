package conversation

import "net/http"

//会话
type Conversation struct {
	C []*http.Cookie
	Client *http.Client
}
