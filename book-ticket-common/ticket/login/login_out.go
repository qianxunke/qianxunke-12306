package login

import (
	"gitee.com/qianxunke/book-ticket-common/basic/utils/conversation"
	"gitee.com/qianxunke/book-ticket-common/proto/user"
	"gitee.com/qianxunke/book-ticket-common/ticket/check_code"
	"log"
	"net/http"
)

//面向客户端的接口
/**
 * 用户注册平台，平台验证，并获取其的联系人信息
 */
func UserRegister(u user.UserInf) (isOk, err error) {
	loginResult := &LoginResult{}
	loginResult.Conversat = &conversation.Conversation{}
	loginResult.Conversat.Client = &http.Client{}
	//验证码
	code, err := check_code.CheckCode(loginResult.Conversat)
	if err != nil {
		log.Println("[CheckCode] error :" + err.Error())
		return
	}
	err = checkUser(u, http.MethodPost, code, loginResult)
	if err != nil {
		log.Println("[checkUser] error :" + err.Error())
		return
	}
	if loginResult.CheckUser {
		err = getToken(loginResult)
		if err != nil {
			log.Println("[getToken] error :" + err.Error())
			return
		}
		if loginResult.GetToken {
			err = checkToken(loginResult)
			if err != nil {
				log.Println("[checkToken] error :" + err.Error())
				return
			}
			if loginResult.CheckToken {
				log.Println("[登陆成功] ok :")
				//获取其联系人

			}
		}
	}
	return
}
