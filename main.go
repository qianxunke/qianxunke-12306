package main

import (
	"qianxunke-12306/modules/book"
	"qianxunke-12306/modules/login"
	"qianxunke-12306/modules/query"
	"qianxunke-12306/modules/sysinit"
	"time"
)

func main() {
	u := login.User{UserName: "dh17862709691", Pwd: "736567805", Id: "5225***********610", Name: "王芳平", TelNum: "18334142052", SeatType: "1", RideDate: "2019-12-15", Departure: "杭州", Terminus: "上海"}
	sysinit.Init()
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
	tran := query.QueryTrainMessage(u, result.Conversat)
	for true {
		if book.Book(*result.Conversat, tran, u) {
			break
		}
		time.Sleep(time.Second * 3)
	}

}
