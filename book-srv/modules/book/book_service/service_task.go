package book_service

import (
	"book-srv/m_client"
	"book-srv/modules/book/book_dao"
	"book-srv/stations"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/basic/utils/conversation"
	"gitee.com/qianxunke/book-ticket-common/basic/utils/http_util"
	"gitee.com/qianxunke/book-ticket-common/notice/sms"
	"gitee.com/qianxunke/book-ticket-common/proto/task"
	"gitee.com/qianxunke/book-ticket-common/proto/user"
	"gitee.com/qianxunke/book-ticket-common/ticket/book/boo_core"
	"gitee.com/qianxunke/book-ticket-common/ticket/login"
	"gitee.com/qianxunke/book-ticket-common/ticket/query"
	"gitee.com/qianxunke/book-ticket-common/ticket/query/bean"
	"gitee.com/qianxunke/book-ticket-common/ticket/static/api"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	taskChnnel chan task.TaskDetails
	taskOpen   chan bool
)

//开始抢票任务
func (s *service) StartBathTicket() {
	if taskChnnel == nil {
		taskChnnel = make(chan task.TaskDetails, 50)
	}
	defer func() {
		if re := recover(); re != nil {
			log.Println(time.Now().Format("2006-01-02 15:04:05") + " 抢票任务重启")
			s.StartBathTicket()
		}
	}()
	ticker := time.NewTicker(time.Second * 10)
	go DoneQueryTicket()
	go func() {
		for range ticker.C {
			log.Println("定时任务：" + time.Now().Format("2006-01-02 15:04:05"))
			rsp, err := s.GetNeedTicketList(100, 1, 1)
			if err != nil {
				log.Println("定时任务：err" + err.Error())
				continue
			}
			if len(rsp) > 0 {
				for _, task := range rsp {
					taskChnnel <- task
				}
			}
		}
		taskOpen <- true
	}()
	<-taskOpen
}

func DoneQueryTicket() {
	for true {
		ta := <-taskChnnel
		go DoneGo(ta)
	}
}

func DoneGo(ta task.TaskDetails) (err error) {
	d, err := book_dao.GetDao()
	if err != nil {
		err = errors.New(fmt.Sprintf("[DoneGo] error %s", err.Error()))
		return
	}
	lastTask, err := d.FindById(ta.Task.TaskId)
	if err != nil {
		err = errors.New(fmt.Sprintf("[DoneGo] error %s\n", err.Error()))
		return
	}
	defer func() {
		if re := recover(); re != nil {
			if err == nil {
				err = errors.New(fmt.Sprintf("订单号 ： %s, error %s\n", ta.Task.TaskId, err.Error()))
				return
			}
		}
		if err != nil {
			if err.Error() != "非待抢" {
				log.Printf("[DoneGo] error %s\n", err.Error())
				//	ta.Task.Status = 1
				//	err = d.Update(lastTask)
			}
		}
	}()

	if lastTask == nil {
		err = errors.New(fmt.Sprintf("[DoneGo] error %s", "订单信息为空"))
		return
	}
	if lastTask.Task.Status != 1 {
		err = errors.New("非待抢")
		return
	}
	//修改状态
	lastTask.Task.Status = 2
	err = d.Update(lastTask)
	if err != nil {
		err = errors.New(fmt.Sprintf("订单号 ： %s, error %s\n", lastTask.Task.TaskId, err.Error()))
		return
	}
	conversation2 := &conversation.Conversation{}
	conversation2.Client = &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, "https://kyfw.12306.cn/otn/leftTicket/init", nil)
	http_util.SetReqHeader(req)
	rsp, err := conversation2.Client.Do(req)
	if err != nil {
		log.Printf("[QueryTrainMessage] error %v", err)
	}
	_, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("[QueryTrainMessage] error %v", err)
	}
	defer rsp.Body.Close()
	http_util.CookieChange(conversation2, rsp.Cookies())

	//开始循环查询票
	errNum := 0
	var chooiseTran bean.Train
	isGetTran := false
	for true {
		if errNum >= 10 {
			err = errors.New(fmt.Sprintf("订单号 ： %s,查询连续出错10次，停止查询\n", lastTask.Task.TaskId))
			return
		}
		trans, err := queryTrainMessage(conversation2, ta.Task.TrainDates, stations.GetStationValueByKey(ta.Task.FindFrom), stations.GetStationValueByKey(ta.Task.FindTo), ta.Task.Type)
		if err != nil {
			errNum++
			continue
		}
		if len(trans) > 0 {
			//判断是否有合适的票
			for _, item := range trans {
				log.Println("TrainCode :" + item.TrainCode + "  CanBuy: " + item.CanBuy)
				if item.Num == ta.Task.Trips && item.CanBuy == "Y" {
					chooiseTran = item
					isGetTran = true
					break
				}
			}
		}
		if isGetTran {
			break
		}
		time.Sleep(time.Second * 3)
	}
	//登陆
	//获取用户信息
	in := &user.InGetUserInfo{UserId: ta.Task.UserId}
	out, err := m_client.UserClient.GetUserInfo(context.TODO(), in)
	if err != nil {
		log.Printf("获取用户信息出错  ： %s,\n", err.Error())
		return
	}
	//重试三次
	loginErrNum := 0
	var loginResult *login.LoginResult
	for true {
		if loginErrNum > 5 {
			err = errors.New("登陆失败5次，取消抢票")
			return err
		}
		_tmp, err := login.LoginAndCheckToken(*out.UserInf)
		if err != nil {
			loginErrNum++
		} else {
			loginResult = _tmp
			break
		}
		time.Sleep(time.Second * 3)
	}
	//开始抢票
	bookErrNum := 0
	isOk := false
	for true {
		if bookErrNum > 5 {
			err = errors.New("抢票失败5次，取消抢票")
			return err
		}
		if !boo_core.Book(*loginResult.Conversat, chooiseTran, ta) {
			loginErrNum++
		} else {
			isOk = true
			break
		}
		time.Sleep(time.Second * 1)
	}
	if isOk {
		log.Println("抢票成功")
		sms.SendTicketSuccessInfoToUser(out.UserInf.MobilePhone, out.UserInf.UserName)
		lastTask.Task.Status = 3
		//lastTask.Task.OkNo=
		err = d.Update(lastTask)
		if err != nil {
			log.Printf("抢票成功 订单号 ： %s, error %s\n", lastTask.Task.TaskId, err.Error())
			return nil
		}
	}
	return
}

func queryTrainMessage(con *conversation.Conversation, TrainDate string, FindFrom string, FindTo string, Type string) (tran []bean.Train, err error) {
	//ADULT
	req1, _ := http.NewRequest(http.MethodGet, api.Query+"A"+"?leftTicketDTO.train_date="+TrainDate+"&leftTicketDTO.from_station="+FindFrom+"&leftTicketDTO.to_station="+FindTo+"&purpose_codes="+Type, nil)
	http_util.AddReqCookie(con.C, req1)
	http_util.SetReqHeader(req1)
	rsp1, err := con.Client.Do(req1)
	if err != nil {
		log.Printf("[QueryTrainMessage] error %v", err)
		return tran, err
	}
	str, err := ioutil.ReadAll(rsp1.Body)
	if err != nil {
		log.Printf("[QueryTrainMessage] error %v", err)
		return tran, err
	}
	defer rsp1.Body.Close()
	//	log.Printf("[QueryTrainMessage]  %s", string(str))
	if rsp1.StatusCode == http.StatusOK {
		queryItem := &bean.QueryItem{}
		err = json.Unmarshal(str, &queryItem)
		if err != nil {
			log.Printf("[QueryTrainMessage] error %v", err)
			return tran, err
		}
		if len(queryItem.Data.Result) > 0 {
			trans, err := query.FormatQueryMessage(queryItem.Data.Result)
			if err != nil {
				log.Printf("[QueryTrainMessage] error %v", err)
				return tran, err
			}
			query.PrintTrainMessage(trans[1:3])
			return trans, nil
		}
	} else {
		log.Printf("[QueryTrainMessage] net error %d :, %v", rsp1.StatusCode, err)
		err = errors.New("[QueryTrainMessage] net error")
		return tran, err
	}
	err = errors.New("[QueryTrainMessage] net error")
	return tran, err

}
