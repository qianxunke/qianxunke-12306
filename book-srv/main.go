package main

import (
	"book-srv/modules/login"
	"book-srv/modules/query"
	"book-srv/modules/sysinit"
	"gitee.com/qianxunke/book-ticket-common/basic/utils/conversation"
	"gitee.com/qianxunke/book-ticket-common/basic/utils/http_util"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	u := login.User{UserName: "dh17862709691", Pwd: "736567805", Id: "5225***********610", Name: "王芳平", TelNum: "18334142052", SeatType: "1", RideDate: "2019-12-15", Departure: "杭州", Terminus: "上海"}
	sysinit.Init()
	/*
		var result login.LoginResult
		for true {
			res, err := login.Login(u)
			//	result2 := login.Login(login.User{UserName: "laidanchao", Pwd: "lai19920127"})
			if err != nil {
				//	login.LoginOut(res)
				time.Sleep(time.Second * 10)
			} else {
				result = *res
				break
			}
		}
		var tran bean.Train
		for true {
			t, err := query.QueryTrainMessage(u, result.Conversat)
			if err == nil {
				tran = t
				break
			}
		}

		for true {
			if book.Book(*result.Conversat, tran, u) {
				break
			}
			time.Sleep(time.Second * 3)
		}

	*/

	//
	conversation := &conversation.Conversation{}
	conversation.Client = &http.Client{}
	//	check_code.CheckCode(conversation)
	//	query.QueryTrainMessage(u,conversation)
	req, _ := http.NewRequest(http.MethodGet, "https://kyfw.12306.cn/otn/leftTicket/init", nil)
	//http_util.AddReqCookie(conversation.C, req)
	http_util.SetReqHeader(req)
	rsp, err := conversation.Client.Do(req)
	if err != nil {
		log.Printf("[QueryTrainMessage] error %v", err)
	}
	_, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("[QueryTrainMessage] error %v", err)
	}
	defer rsp.Body.Close()
	log.Printf("[QueryTrainMessage]  %v\n", rsp.Cookies())
	//log.Printf("[QueryTrainMessage]  %s\n", string(str))
	http_util.CookieChange(conversation, rsp.Cookies())
	query.QueryTrainMessage(u, conversation)
}
