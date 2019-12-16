package login

import (
	"encoding/json"
	"errors"
	"gitee.com/qianxunke/book-ticket-common/basic/utils/conversation"
	"gitee.com/qianxunke/book-ticket-common/basic/utils/http_util"
	"gitee.com/qianxunke/book-ticket-common/proto/user"
	bookBean "gitee.com/qianxunke/book-ticket-common/ticket/book/bean"
	"gitee.com/qianxunke/book-ticket-common/ticket/check_code"
	"gitee.com/qianxunke/book-ticket-common/ticket/static/api"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type LoginResult struct {
	Conversat  *conversation.Conversation
	CheckUser  bool   // 用户是否可以登录
	GetToken   bool   // 是否获取到Token
	CheckToken bool   // Token是否通过检查
	Login      bool   // 是否完成登录
	Newapptk   string // 检查Token发送的识别码
	Apptk      string // 登陆成功的识别码
	Username   string // 登陆成功后获取的用户名

}

func checkUser(u user.UserInf, method string, checkCode string, loginResult *LoginResult) (err error) {
	// 准备URL
	log.Println("正在验证账号密码...")
	data := url.Values{}
	data.Set("username", u.TranUserAccount)
	data.Set("password", u.TranUserPwd)
	data.Set("appid", "otn")
	data.Set("answer", checkCode)
	req, _ := http.NewRequest(method, api.LoginURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	req.Header.Set("Referer", "https://kyfw.12306.cn/otn/resources/login.html")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	http_util.SetReqHeader(req)
	http_util.AddReqCookie(loginResult.Conversat.C, req)
	rsp, err := loginResult.Conversat.Client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if rsp.StatusCode == http.StatusOK {
		str, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		m := make(map[string]interface{})
		err = json.Unmarshal(str, &m)
		if err != nil {
			if method == http.MethodPut {
				return err
			} else {
				log.Println("返回了html码，尝试put请求")
				return checkUser(u, http.MethodPut, checkCode, loginResult)
			}
		}
		log.Printf("[checkUser]%v", m)
		isOk := false
		switch m["result_code"].(type) {
		case string:
			if m["result_code"].(string) == "0" {
				isOk = true
			}
		case float64:
			if m["result_code"].(float64) == 0 {
				isOk = true
			}
		}
		if isOk {
			log.Println("账号密码正确...")
			loginResult.CheckUser = true
			http_util.CookieChange(loginResult.Conversat, rsp.Cookies())
			return err
		} else {
			log.Println("账号密码错误...")
			loginResult.CheckUser = false
			http_util.CookieChange(loginResult.Conversat, rsp.Cookies())
			return err
		}
	} else {
		log.Println("连接错误...")
		return err
	}
}

/**
 * 获取登录Token
 *
 * @param client
 * @return
 */
func getToken(loginResult *LoginResult) (err error) {
	// 准备URL
	log.Println("正在获取Token...")
	data := url.Values{}
	data.Set("appid", "otn")
	req, _ := http.NewRequest(http.MethodPost, api.GetToken, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	req.Header.Set("Referer", "https://kyfw.12306.cn/otn/resources/login.html")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	http_util.AddReqCookie(loginResult.Conversat.C, req)
	http_util.SetReqHeader(req)
	rsp, err := loginResult.Conversat.Client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if rsp.StatusCode == http.StatusOK {
		str, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		m := make(map[string]interface{})
		err = json.Unmarshal(str, &m)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		log.Println(m)
		http_util.CookieChange(loginResult.Conversat, rsp.Cookies())
		loginResult.Newapptk = m["newapptk"].(string)
		loginResult.GetToken = true
		return err
	} else {
		str, err := ioutil.ReadAll(rsp.Body)
		err = errors.New(string(str))
		return err
	}

}

/**
 * 登录
 * @param loginResult
 * @return
 */

func checkToken(loginResult *LoginResult) (err error) {
	log.Println("正在验证Token...")
	data := url.Values{}
	data.Set("tk", loginResult.Newapptk)
	req, _ := http.NewRequest(http.MethodPost, api.CheckToken, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	req.Header.Set("Referer", "https://kyfw.12306.cn/otn/passport?redirect=/otn/login/userLogin")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	http_util.SetReqHeader(req)
	http_util.AddReqCookie(loginResult.Conversat.C, req)
	log.Println()
	rsp, err := loginResult.Conversat.Client.Do(req)
	if err != nil {
		return err
	}
	str, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Println("[checkToken] error" + err.Error())
		return err
	}
	log.Printf("[checkToken] respose body %s\n", string(str))
	if rsp.StatusCode == http.StatusOK {
		m := make(map[string]interface{})
		err = json.Unmarshal(str, &m)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		isOk := false
		switch m["result_code"].(type) {
		case string:
			if m["result_code"].(string) == "0" {
				isOk = true
			}
		case float64:
			if m["result_code"].(float64) == 0 {
				isOk = true
			}
		}
		if isOk {
			loginResult.CheckToken = true
			loginResult.Apptk = m["apptk"].(string)
			loginResult.Username = m["username"].(string)
			http_util.CookieChange(loginResult.Conversat, rsp.Cookies())
			log.Println("登陆成功,用户名:" + loginResult.Username)
			return err
		} else {
			log.Println("请求失败...,需要重新登陆")
			return errors.New("需要重新登陆")
		}
	} else {
		log.Printf("网络错误，错误信息: %d", rsp.StatusCode)
		err = errors.New("网络错误，错误信息:")
		return err
	}
}

//检查用户登陆状态

func CheckUserStatus(bookResult *bookBean.BookResult, conversation *conversation.Conversation) (err error) {
	defer func() {
		if re := recover(); re != nil {
			log.Printf("[CheckUserStatus]: %v", re)
			err = errors.New("CheckUserStatus error")

		}
	}()
	log.Println("正在检查用户是否登录...")
	data := url.Values{}
	data.Set("_json_att", "")
	req, _ := http.NewRequest(http.MethodPost, api.CheckLoginStatus, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	req.Header.Set("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	http_util.SetReqHeader(req)
	http_util.AddReqCookie(conversation.C, req)
	rsp, err := conversation.Client.Do(req)
	if err != nil {
		log.Printf("[CheckUserStatus]: %v", err.Error())
		return err
	}
	if rsp.StatusCode == http.StatusOK {
		http_util.CookieChange(conversation, rsp.Cookies())
		bodyBytes, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Printf("[CheckUserStatus]: %v", err.Error())
			return err
		}
		defer rsp.Body.Close()
		log.Printf("[CheckUserStatus]: response : %v", string(bodyBytes))
		m := make(map[string]interface{})
		err = json.Unmarshal(bodyBytes, &m)
		if err != nil {
			log.Printf("[CheckUserStatus]: %v", err.Error())
			return err
		}
		if m["status"].(bool) && m["data"].(map[string]interface{})["flag"].(bool) {
			bookResult.CheckUser = true
		} else {
			return errors.New("login status error")
		}
	} else {
		log.Printf("[CheckUserStatus]: response error %d", rsp.StatusCode)
		return errors.New("response error")
	}
	return err
}

func LoginOut(loginResult *LoginResult) (err error) {
	log.Println("正在退出登陆...")
	req, _ := http.NewRequest(http.MethodGet, api.LoginOut, nil)
	http_util.SetReqHeader(req)
	http_util.AddReqCookie(loginResult.Conversat.C, req)
	rsp, err := loginResult.Conversat.Client.Do(req)
	if err != nil {
		log.Printf("[CheckUserStatus]: %v", err.Error())
		return err
	}
	defer rsp.Body.Close()
	req, _ = http.NewRequest(http.MethodGet, api.Init, nil)
	http_util.SetReqHeader(req)
	http_util.AddReqCookie(loginResult.Conversat.C, req)
	rsp, err = loginResult.Conversat.Client.Do(req)
	defer rsp.Body.Close()
	if err != nil {
		log.Printf("[LoginOut]: %v", err.Error())
		return err
	}
	loginResult.Conversat.C = rsp.Cookies()
	data := &url.Values{}
	data.Set("appid", "otn")
	req, _ = http.NewRequest(http.MethodPost, "https://kyfw.12306.cn/passport/web/auth/uamtk", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	http_util.SetReqHeader(req)
	http_util.AddReqCookie(loginResult.Conversat.C, req)
	rsp, err = loginResult.Conversat.Client.Do(req)
	defer rsp.Body.Close()
	if err != nil {
		log.Printf("[LoginOut]: %v", err.Error())
		return err
	}
	loginResult.Conversat.Client.CloseIdleConnections()
	return
}

/**
  登陆并检查token
*/
func LoginAndCheckToken(u user.UserInf) (loginResult *LoginResult, err error) {
	loginResult = &LoginResult{}
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
				return
			}
		}
	}
	return
}
