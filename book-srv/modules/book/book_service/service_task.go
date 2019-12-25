package book_service

import (
	"book-srv/m_client"
	"book-srv/modules/book/book_dao"
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
	"net/url"
	"strings"
	"time"
)

var (
	taskChnnel chan task.Task
	taskOpen   chan bool

	errorChnnel chan task.Task
	errorOpen   chan bool
)

//1.时间过了就不用抢
//2.还不放票的日期不能抢 （需要弄个缓存，redis
//3.多个时间需要循环判断一下，
//4.多个车次循环判断
//5.6点前，10点后不能抢
//6，查询地址经常变动  （redis控制）
//，异常退出的任务需要重新进，状态为抢票中，但刷新时间过了一个小时

//再造一个定时器，处理异常任务
func (s *service) StartBathDoneError() {
	if errorChnnel == nil {
		errorChnnel = make(chan task.Task, 50)
	}
	defer func() {
		if re := recover(); re != nil {
			log.Println(time.Now().Format("2006-01-02 15:04:05") + " 抢票任务重启")
			s.StartBathDoneError()
		}
	}()
	ticker := time.NewTicker(time.Second * 60 * 5)
	go func() {
		for true {
			ta := <-errorChnnel
			go DoneErrorTask(ta)
		}
	}()
	go func() {
		for range ticker.C {
			log.Println("定时处理异常任务-------->：" + time.Now().Format("2006-01-02 15:04:05"))
			//如果非抢票时间就跳过
			d, err := book_dao.GetDao()
			if err != nil {
				log.Println(err.Error())
				return
			}
			rsp, err := d.ExceptionQuery(100, 1)
			if err != nil {
				log.Println("定时任务：err" + err.Error())
				continue
			}
			//	log.Printf("定时处理异常任务,异常包裹：%v\n", rsp)
			if len(rsp) > 0 {
				for _, ta := range rsp {
					errorChnnel <- ta
				}
			}
		}
		errorOpen <- true
	}()
	<-errorOpen
}

func DoneErrorTask(ta task.Task) {
	//3已完成  0取消 2抢票中 1待抢票，4再等等（针对预约时间还有很长时间的
	//抢票中 如果半小时还没更新时间并且还没过抢票时间，就转化为待抢票
	//把再等等更新为待抢票
	d, err := book_dao.GetDao()
	if err != nil {
		log.Println(err.Error())
		return
	}
	lastTask, err := d.GetTask(ta.TaskId)
	if err != nil {
		log.Println(err.Error())
		return
	}

	//如果是抢票中，且半小时还没更新，
	if lastTask.Status == 2 && (time.Now().Unix()-lastTask.UpdateTime > 300) {
		//判断抢票日期是否已经失效
		s1 := strings.Split(lastTask.TrainDates, ",")
		var ok bool
		for _, item := range s1 {
			//如果还可以抢票
			if isCanQuery(item) > 0 {
				ok = true
				break
			}
		}
		if ok {
			//如果满足条件则更新任务，重新回到抢票中
			_ = d.UpdateStatus(lastTask.GetTaskId(), 1)
			return
		}
	}
	//如果是还没到抢票日期，再次计算是否可以抢票啦
	if lastTask.Status == 4 {
		s1 := strings.Split(lastTask.TrainDates, ",")
		var ok bool
		for _, item := range s1 {
			days := isCanQuery(item)
			//如果还可以抢票
			if days > 0 && days < 30 {
				ok = true
				break
			}
		}
		if ok {
			//如果满足条件则更新任务，重新回到抢票中
			_ = d.UpdateStatus(lastTask.GetTaskId(), 1)
			return
		}
	}

}

//开始抢票任务
func (s *service) StartBathTicket() {
	if taskChnnel == nil {
		taskChnnel = make(chan task.Task, 20)
	}
	defer func() {
		if re := recover(); re != nil {
			log.Println(time.Now().Format("2006-01-02 15:04:05") + " 抢票任务重启")
			s.StartBathTicket()
		}
	}()
	ticker := time.NewTicker(time.Second * 60)
	go DoneQueryTicket()
	go func() {
		for range ticker.C {
			log.Println("定时任务：" + time.Now().Format("2006-01-02 15:04:05"))
			//如果非抢票时间就跳过
			if !timeIsOk() {
				continue
			}
			rsp, err := s.GetNeedTicketList(10, 1, 1)
			if err != nil {
				log.Println("定时任务：err" + err.Error())
				continue
			}
			if len(rsp) > 0 {
				for _, ta := range rsp {
					taskChnnel <- ta
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

func DoneGo(ta task.Task) (err error) {
	d, err := book_dao.GetDao()
	if err != nil {
		err = errors.New(fmt.Sprintf("[DoneGo] error %s", err.Error()))
		return
	}
	lastTask, err := d.FindById(ta.TaskId)
	if err != nil {
		err = errors.New(fmt.Sprintf("[DoneGo] error %s\n", err.Error()))
		return
	}
	defer func() {
		if re := recover(); re != nil {
			if err == nil {
				err = errors.New(fmt.Sprintf("订单号 ： %s, recover %s\n", ta.TaskId, err.Error()))
				return
			}
		}
	}()

	if lastTask == nil {
		err = errors.New(fmt.Sprintf("[DoneGo] error %s", "订单信息为空"))
		return
	}
	//如果为0失效就不需要抢
	if lastTask.Task.Status == 0 {
		return
	}
	//修改状态 正在抢票
	//为了安全，再判断当前时间是否可抢 5,取消啦
	if lastTask.Task.Status == 5 {
		return
	}
	//如果为 2已经有人在抢啦，就不需要抢
	if lastTask.Task.Status == 2 {
		return
	}
	//如果等于4，判断是否可以抢了
	if lastTask.Task.Status == 4 {
		var ok bool
		dTemp := strings.Split(lastTask.Task.TrainDates, ",")
		for _, item := range dTemp {
			d := isCanQuery(item)
			if d > 0 && d < 30 {
				ok = true
			}
		}
		if !ok {
			return
		}
	}
	lastTask.Task.Status = 2
	err = d.UpdateStatus(lastTask.Task.GetTaskId(), 2)
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
		_ = d.UpdateStatus(lastTask.Task.GetTaskId(), 1)
		log.Printf("[QueryTrainMessage] error %v", err)
	}
	if rsp.StatusCode != http.StatusOK {
		log.Printf("[https://kyfw.12306.cn/otn/leftTicket/init] error %d\n", rsp.StatusCode)
		//先让该单锁5分钟吧
		return
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		_ = d.UpdateStatus(lastTask.Task.GetTaskId(), 1)
		log.Printf("[QueryTrainMessage] error %v", err)
	}
	defer rsp.Body.Close()
	str := string(body)
	var CLeftTicketUrl string
	htmls := strings.Split(str, "\n")
	for _, line := range htmls {
		if strings.Contains(line, "CLeftTicketUrl") {
			CLeftTicketUrl = line[strings.Index(line, "'")+1 : len(line)-2]
			break
		}
	}
	if len(CLeftTicketUrl) == 0 {
		log.Println("CLeftTicketUrl : 空")
		return
	}

	http_util.CookieChange(conversation2, rsp.Cookies())
	//开始循环查询票
	errNum := 0
	var chooiseTran bean.Train
	var days []string
	//获取可用日期
	dayTmp := strings.Split(lastTask.Task.TrainDates, ",")
	for _, item := range dayTmp {
		d := isCanQuery(item)
		if d > 0 && d < 30 {
			days = append(days, item)
		}
	}
	if len(days) == 0 {
		log.Println("没有可用日期了，直接设置失败")
		_ = d.UpdateStatus(lastTask.Task.GetTaskId(), 0)
		return
	}
	//日期循环
	//var d_index int
	//获取用户信息
	in := &user.InGetUserInfo{UserId: lastTask.Task.UserId}
	out, err := m_client.UserClient.GetUserInfo(context.TODO(), in)
	if err != nil {
		_ = d.UpdateStatus(lastTask.Task.GetTaskId(), 1)
		log.Printf("获取用户信息出错  ： %s,\n", err.Error())
		return
	}
	isOk := false
	for true {
		if errNum >= 10 {
			log.Printf("用户 ：%s ,订单号 ： %s,查询连续出错10次，停止查询\n", out.UserInf.MobilePhone, lastTask.Task.TaskId)
			err = d.UpdateStatus(lastTask.Task.GetTaskId(), 1)
			return err
		}
		//判断该订单是否已经取消
		t1, err := d.GetTask(ta.TaskId)
		if err != nil {
			return err
		}
		if t1.Status != 2 {
			return errors.New("[取消抢票]")
		}
		err = d.UpdateStatus(lastTask.Task.GetTaskId(), 2)
		if err != nil {
			return err
		}
		if !timeIsOk() {
			err = d.UpdateStatus(lastTask.Task.GetTaskId(), 1)
			return err
		}
		for i := 0; i < len(days); i++ {

			trans, err := queryTrainMessage(CLeftTicketUrl, conversation2, days[i], lastTask.Task.FindFrom, lastTask.Task.FindTo, lastTask.Task.Type)
			if err != nil {
				errNum++
				continue
			}
			if len(trans) > 0 {
				//判断是否有合适的票
				for _, item := range trans {
					//log.Println("TrainCode :" + item.Num + "  CanBuy: " + item.CanBuy)
					if strings.Contains(lastTask.Task.Trips, item.Num) && item.CanBuy == "Y" {
						//判断是否有票，
						chooiseTran = item
						/**
						PASSENGER_TICKER_STR = {
						'一等座': 'M',
						 '特等座': 'P',
						'二等座': 'O',
						'商务座': 9,
						'硬座': 1,
						'无座': 1,
						 '软座': 2,
						 '软卧': 4,
						 硬卧': 3,
						}
						*/
						var ishaveSet bool
						if strings.Contains(item.Num, "D") || strings.Contains(item.Num, "G") {
							if item.Edz != "" && item.Edz != "无" {
								if strings.Contains(lastTask.Task.SeatTypes, "O") {
									ishaveSet = true
								}
							}
							if item.Ydz != "" && item.Ydz != "无" {
								if strings.Contains(lastTask.Task.SeatTypes, "M") {
									ishaveSet = true
								}
							}
							if item.Swtdz != "" && item.Swtdz != "无" {
								if strings.Contains(lastTask.Task.SeatTypes, "9") {
									ishaveSet = true
								}
							}

						} else {
							if item.Yz != "" && item.Yz != "无" {
								if strings.Contains(lastTask.Task.SeatTypes, "1") {
									ishaveSet = true
								}
							}
							if item.Wz != "" && item.Wz != "无" {
								if strings.Contains(lastTask.Task.SeatTypes, "S") {
									ishaveSet = true
								}
							}
							if item.Yw != "" && item.Yw != "无" {
								if strings.Contains(lastTask.Task.SeatTypes, "3") {
									ishaveSet = true
								}
							}
							if item.Rw != "" && item.Rw != "无" {
								if strings.Contains(lastTask.Task.SeatTypes, "4") {
									ishaveSet = true
								}
							}
							if item.Rz != "" && item.Rz != "无" {
								if strings.Contains(lastTask.Task.SeatTypes, "2") {
									ishaveSet = true
								}
							}
						}
						if ishaveSet {
							//登陆
							//重试三次
							var loginResult *login.LoginResult
							for i := 3; i >= 0; i-- {
								_tmp, err := login.LoginAndCheckToken(*out.UserInf)
								if err == nil {
									loginResult = _tmp
									break
								}
								time.Sleep(time.Second * 3)
							}
							if loginResult == nil {
								break
							}
							//开始抢票
							bookErrNum := 0
							for true {
								if bookErrNum > 2 {
									err = d.UpdateStatus(lastTask.Task.GetTaskId(), 1)
									err = errors.New("抢票失败3次，取消抢票,重新入队列抢票")
									return err
								}
								if !boo_core.Book(*loginResult.Conversat, chooiseTran, *lastTask) {
									bookErrNum++
								} else {
									isOk = true
									break
								}
								time.Sleep(time.Second * 3)
							}
							if isOk {
								break
							}
						}
					}
				}
			}
			if isOk {
				break
			}
			//判断当前时间是否还有效
			d := isCanQuery(days[i])
			if d < 0 || d > 30 {
				days = append(days[:i], days[(i+1):]...)
			}
			time.Sleep(time.Second * 5)
		}

		if isOk {
			break
		}

		if len(days) == 0 {
			break
		}
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
	} else {
		err = d.UpdateStatus(lastTask.Task.GetTaskId(), 1)
	}
	return
}

func queryTrainMessage(CLeftTicketUrl string, con *conversation.Conversation, TrainDate string, FindFrom string, FindTo string, Type string) (tran []bean.Train, err error) {
	//ADULT
	defer func() {
		if re := recover(); re != nil {
			if err == nil {
				err = errors.New("[queryTrainMessage] : 查询异常")
				return
			}
		}
	}()
	s, _ := url.PathUnescape(api.Query + CLeftTicketUrl + "?leftTicketDTO.train_date=" + TrainDate + "&leftTicketDTO.from_station=" + FindFrom + "&leftTicketDTO.to_station=" + FindTo + "&purpose_codes=" + Type)
	req1, _ := http.NewRequest(http.MethodGet, s, nil)
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

//当前是否可以抢票
func timeIsOk() (ok bool) {
	now := time.Now()
	if now.Hour() >= 6 && now.Hour() < 23 {
		return true
	}
	return false
}

//判断已经可以买票，可以提前30天
func isCanQuery(trainDate string) float64 {
	a, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	b, _ := time.Parse("2006-01-02", trainDate)
	d := b.Sub(a)
	return d.Hours() / 24

}
