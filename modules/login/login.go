package login

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"qianxunke-12306/common/http_util"
	"qianxunke-12306/config/api"
	"qianxunke-12306/modules/check_code"
	"qianxunke-12306/modules/conversation"
	"strconv"
	"strings"
)

type User struct {
	UserName  string // 12306用户名
	Pwd       string // 12306密码
	Name      string // 乘车人用户名
	Id        string // 乘车人身份证
	TelNum    string // 语音通知及接收短信手机号
	SeatType  string // 席别
	SeatNum   string // 座号，同12306，不一定可以选到希望的座位
	RideDate  string // 乘车日期 格式为2017.01.01
	Departure string // 查询始发站
	Terminus  string // 查询终点站
	StartTime string // 最早乘车时间
	EndTime   string // 最晚乘车时间

}

type LoginResult struct {
	Conversat   *conversation.Conversation
	CheckUser  bool   // 用户是否可以登录
	GetToken   bool   // 是否获取到Token
	CheckToken bool   // Token是否通过检查
	Login      bool   // 是否完成登录
	Newapptk   string // 检查Token发送的识别码
	Apptk      string // 登陆成功的识别码
	Username   string // 登陆成功后获取的用户名

}

func checkUser(u User, checkCode string, loginResult *LoginResult) (err error) {
	// 准备URL
	log.Println("正在验证账号密码...")
	data:=url.Values{}
	data.Set("username",u.UserName)
	data.Set("password",u.Pwd)
	data.Set("appid","otn")
	data.Set("answer",checkCode)
	req, _ := http.NewRequest(http.MethodPut,api.LoginURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type","application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept","application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Origin","https://kyfw.12306.cn")
	req.Header.Set("Referer","https://kyfw.12306.cn/otn/resources/login.html")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	http_util.SetReqHeader(req)
	http_util.AddReqCookie(loginResult.Conversat.C, req)

	log.Printf("req: %+v",req)

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
		log.Printf("rsp: %+v",string(str))
		err = json.Unmarshal(str, &m)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		log.Println(m)
		if m["result_code"] == 0 {
			log.Println("账号密码正确...")
			loginResult.CheckUser = true
			loginResult.Conversat.C = rsp.Cookies()
			return err
		} else {
			log.Println("账号密码正确...")
			loginResult.CheckUser = false
			loginResult.Conversat.C = rsp.Cookies()
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

	req, _ := http.NewRequest(http.MethodPost, api.GetToken, nil)
	req.Form=url.Values{}
	req.Form.Set("appid", "otn")
	http_util.AddReqCookie(loginResult.Conversat.C, req)
	http_util.SetReqHeader(req)
	rsp, err := loginResult.Conversat.Client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if rsp.StatusCode == http.StatusOK {
		log.Println("获取Token成功...");
		str, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		m := make(map[string]string)
		err = json.Unmarshal(str, &m)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		loginResult.Newapptk = m["newapptk"]
		loginResult.GetToken = true
		return err
	} else {
		log.Println("网络错误，错误信息:")
		err = errors.New("网络错误，错误信息:")
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
	req, _ := http.NewRequest(http.MethodPost, api.CheckToken, nil)
	req.Form=url.Values{}
	req.Form.Set("tk", loginResult.Newapptk)
	http_util.SetReqHeader(req)
	http_util.AddReqCookie(loginResult.Conversat.C, req)

	rsp, err := loginResult.Conversat.Client.Do(req)
	if err != nil {
		return err
	}
	if rsp.StatusCode == http.StatusOK {
		str, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Println("[checkToken] error"+err.Error())
			return err
		}
		m := make(map[string]string)
		err = json.Unmarshal(str, &m)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		if m["result_code"] == "0" {
			loginResult.CheckToken = true
			loginResult.Apptk = m["apptk"]
			loginResult.Username = m["username"]
			loginResult.Conversat.C = rsp.Cookies()
			log.Println("登陆成功,用户名:" + loginResult.Username)
			return err
		} else {
			log.Println("登陆失败,用户名:" + loginResult.Username)
			return err
		}
	} else {
		log.Println("网络错误，错误信息:")
		err = errors.New("网络错误，错误信息:")
		return err
	}
}

func Login()  {
	loginResult:=&LoginResult{}
	loginResult.Conversat=&conversation.Conversation{}
	loginResult.Conversat.Client=&http.Client{}

	//验证码
	code,err:=check_code.CheckCode(loginResult.Conversat)
	if err!=nil{
		log.Println("[CheckCode] error :"+err.Error())
		return
	}
	log.Println("验证码 ok :"+code)
    err=checkUser(User{UserName:"dh17862709691",Pwd:"736567805"},code,loginResult)
    if err!=nil{
		log.Println("[checkUser] error :"+err.Error())
		return
	}
	if loginResult.CheckUser {
		err=getToken(loginResult)
		if err!=nil{
			log.Println("[getToken] error :"+err.Error())
			return
		}
		if loginResult.GetToken {
			err=checkToken(loginResult)
			if err!=nil{
				log.Println("[checkToken] error :"+err.Error())
				return
			}
			if loginResult.CheckToken{
				log.Println("[登陆成功] ok :")
				return
			}
		}
	}
}


