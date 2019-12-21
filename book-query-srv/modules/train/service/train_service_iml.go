package service

import (
	trainDao "book-query-srv/modules/train/dao"
	"encoding/json"
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/basic/utils/conversation"
	"gitee.com/qianxunke/book-ticket-common/basic/utils/http_util"
	ticketProto "gitee.com/qianxunke/book-ticket-common/proto/ticket"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//获取信息
func (s *service) GetTrainById(req *ticketProto.In_GetTrainInfo) (rsp *ticketProto.Out_GetTrainInfo) {
	rsp = &ticketProto.Out_GetTrainInfo{}
	rsp.Error = &ticketProto.Error{}
	if len(req.SecretStr) <= 0 {
		rsp.Error = &ticketProto.Error{
			Code:    http.StatusBadRequest,
			Message: "请求参数有误！",
		}
		return
	}
	dao, err := trainDao.GetTicketDao()
	if err != nil {
		rsp.Error = &ticketProto.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	rsp.Train, err = dao.FindById(req.SecretStr)
	if err != nil {
		rsp.Error = &ticketProto.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	rsp.Error = &ticketProto.Error{
		Code:    http.StatusOK,
		Message: "查询成功！",
	}
	return

}

//修改信息
func (s *service) UpdateTrainInfo(req *ticketProto.In_UpdateTrainInfo) (rsp *ticketProto.Out_UpdateTrainInfo) {
	defer func() {
		if re := recover(); re != nil {
			rsp.Error = &ticketProto.Error{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("[UpdateTrainInfo] error %v", re),
			}
		}
	}()
	rsp = &ticketProto.Out_UpdateTrainInfo{}
	dao, err := trainDao.GetTicketDao()
	if err != nil {
		rsp.Error = &ticketProto.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	err = dao.Update(req.Train)
	if err != nil {
		rsp.Error = &ticketProto.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	rsp.Error = &ticketProto.Error{
		Code:    http.StatusOK,
		Message: "修改成功！",
	}
	return
}

func GetClent() (con *conversation.Conversation) {
	con = &conversation.Conversation{}
	con.Client = &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, "https://kyfw.12306.cn/otn/leftTicket/init", nil)
	http_util.SetReqHeader(req)
	rsp, err := con.Client.Do(req)
	if err != nil {
		log.Printf("[QueryTrainMessage] error %v", err)
	}
	_, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("[QueryTrainMessage] error %v", err)
	}
	defer rsp.Body.Close()
	http_util.CookieChange(con, rsp.Cookies())
	return
}

//获取列表
func (s *service) GetTrainInfoList(req *ticketProto.In_GetTrainInfoList) (rsp *ticketProto.Out_GetTrainInfoList) {
	rsp = &ticketProto.Out_GetTrainInfoList{}
	var err error
	dao, err := trainDao.GetTicketDao()
	if err != nil {
		log.Printf("ERROR: %v", err)
		rsp.Error = &ticketProto.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	value, err := dao.SimpleQuery(req.FindFrom, req.FindTo, req.TrainDate, req.PurposeCodes)
	if err != nil || len(value) == 0 {
		log.Printf("ERROR: %v", err)
		//从12306查询
		q, err := dao.GetRedisClient().Do("GET", "10086").String()
		if err != nil {
			log.Printf("redis : %v\n", err)
			q = "Z"
		}
		time.Sleep(time.Millisecond * 500)
		rsp.TrainList, err = s.queryTrainMessage(q, *req)
		if err != nil {
			log.Printf("ERROR: %v\n", err)
		} else {
		_:
			dao.Insert(rsp.TrainList)
		}
	} else {
		//将json转对象
		_ = json.Unmarshal([]byte(value), rsp.TrainList)

	}
	var message string
	if rsp.Total > 0 && len(rsp.TrainList) > 0 {
		message = "查询成功！"
	} else {
		message = "没有查到任何数据！"
	}
	//统计有多少条
	rsp.Error = &ticketProto.Error{
		Code:    http.StatusOK,
		Message: message,
	}
	rsp.Limit = -1
	rsp.Pages = -1
	return

}

//新建信息
func (s *service) CreateTrain(req *ticketProto.In_AddTrain) (rsp *ticketProto.Out_AddTrain) {
	rsp = &ticketProto.Out_AddTrain{}
	dao, err := trainDao.GetTicketDao()
	if err != nil {
		rsp.Error = &ticketProto.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	err = dao.Insert(nil)
	if err != nil {
		rsp.Error = &ticketProto.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	rsp.Error = &ticketProto.Error{
		Code:    http.StatusOK,
		Message: "新增成功!",
	}
	return
}
