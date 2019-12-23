package boo_core

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/basic/utils/conversation"
	"gitee.com/qianxunke/book-ticket-common/basic/utils/http_util"
	"gitee.com/qianxunke/book-ticket-common/proto/task"
	bookBean "gitee.com/qianxunke/book-ticket-common/ticket/book/bean"
	"gitee.com/qianxunke/book-ticket-common/ticket/static/api"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/**
 * 提交订单请求 检查该用户是否有未完成订单，如果有则返回false
 *
 * @param bookResult
 * @return
 */
func submitOrder(conversation *conversation.Conversation, bookResult *bookBean.BookResult) (err error) {
	defer func() {
		if re := recover(); re != nil {
			err = errors.New(fmt.Sprintf("%v", re))
		}
	}()
	log.Println("正在检查是否有未完成订单...")
	data := &url.Values{}
	//tm := time.Now()
	str, _ := url.QueryUnescape(bookResult.SelectTran.SecretStr)
	data.Set("secretStr", str) //这里注意
	data.Set("train_date", bookResult.SelectTran.StartDate[:4]+"-"+bookResult.SelectTran.StartDate[4:6]+"-"+bookResult.SelectTran.StartDate[6:])
	//	SimpleDateFormat sdf = new SimpleDateFormat("yyyy-MM-dd");
	//	String back_train_date = sdf.format(new Date());

	data.Set("back_train_date", "")
	data.Set("tour_flag", "dc")
	data.Set("purpose_codes", bookResult.Task.Task.Type)
	data.Set("query_from_station_name", bookResult.SelectTran.FindFrom)
	data.Set("query_to_station_name", bookResult.SelectTran.FindTo)
	data.Set("undefined", "")
	req, _ := http.NewRequest(http.MethodPost, api.SubmitOrderRequestURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	req.Header.Set("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	http_util.SetReqHeader(req)
	http_util.AddReqCookie(conversation.C, req)
	log.Printf("[submitOrder]: req %v\n", data)
	rsp, err := conversation.Client.Do(req)
	if err != nil {
		log.Printf("[submitOrder]: %s\n", err.Error())
		return
	}
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("[submitOrder]: %s\n", err.Error())
		return
	}
	defer rsp.Body.Close()
	//	log.Printf("[submitOrder] bodyBytes : %s\n", string(bodyBytes))
	if rsp.StatusCode == http.StatusOK {
		http_util.CookieChange(conversation, rsp.Cookies())
		m := make(map[string]interface{})
		err = json.Unmarshal(bodyBytes, &m)
		if err != nil {
			log.Printf("[submitOrder]: %s\n", err.Error())
			return
		}
		if m["status"].(bool) {
			bookResult.SubmitOrder = true
		} else {
			log.Printf("[submitOrder]: %v\n", m["messages"])
			bookResult.SubmitOrder = false
		}
	} else {
		log.Printf("[submitOrder]:  %s\n", string(bodyBytes))
		return
	}

	return

}

func GetInitDc(conversation *conversation.Conversation) (GlobalRepeatSubmitToken string, da bookBean.InitDC, err error) {
	defer func() {
		if re := recover(); re != nil {
			if err == nil {
				err = errors.New(fmt.Sprintf("[GetInitDc] recover %v", re))
				return
			}
		}
	}()
	log.Println("正在请求InitDc...")
	data := &url.Values{}
	data.Set("_json_att", "")
	req, _ := http.NewRequest(http.MethodPost, api.InitDcURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	req.Header.Set("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	http_util.SetReqHeader(req)
	http_util.AddReqCookie(conversation.C, req)
	rsp, err := conversation.Client.Do(req)
	if err != nil {
		log.Printf("[getInitDc]: %s\n", err.Error())
		return
	}
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("[getInitDc]: %s\n", err.Error())
		return
	}
	defer rsp.Body.Close()
	//	log.Printf("[getInitDc] bodyBytes : %s\n", string(bodyBytes))
	if rsp.StatusCode == http.StatusOK {
		http_util.CookieChange(conversation, rsp.Cookies())
		str := string(bodyBytes)
		htmls := strings.Split(str, "\n")
		for _, line := range htmls {
			if strings.Contains(line, "globalRepeatSubmitToken") {
				//	log.Println("line : "+line )
				GlobalRepeatSubmitToken = line[strings.Index(line, "'")+1 : len(line)-2]
				log.Println("bookResult.GlobalRepeatSubmitToken : " + GlobalRepeatSubmitToken)
			}

			if strings.Contains(line, "var ticketInfoForPassengerForm") {
				//	log.Println("line : "+line )
				ticketInfo := line[strings.Index(line, "{") : len(line)-1]
				//	log.Printf("ticketInfo : %v\n" ,ticketInfo)
				da, err = bookBean.FormatInitDc(ticketInfo)
				if err != nil {
					return
				}
				//	log.Printf("bookResult.InitDcInfo : %v\n" ,bookResult.InitDcInfo)
			}
		}
		//bookResult.InitDc = true
	} else {
		log.Printf("[getInitDc]: error  %s\n", rsp.Status)
		return
	}

	return

}

func GetPassenger(method string, conversation *conversation.Conversation, GlobalRepeatSubmitToken string) (Passenger bool, ps []bookBean.Normal_passengers, err error) {
	defer func() {
		if re := recover(); re != nil {
			if err == nil {
				err = errors.New(fmt.Sprintf("[GetPassenger] %v", re))
				return
			}
		}
	}()
	log.Println("正在请求乘客信息...")
	data := &url.Values{}
	data.Set("_json_att", "")
	data.Set("REPEAT_SUBMIT_TOKEN", GlobalRepeatSubmitToken)
	req, _ := http.NewRequest(method, api.GetPassenger, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	//	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	req.Header.Set("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	http_util.SetReqHeader(req)
	http_util.AddReqCookie(conversation.C, req)
	rsp, err := conversation.Client.Do(req)
	if err != nil {
		log.Printf("[getPassenger]: %s\n", err.Error())
		return
	}
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("[getPassenger]: %s\n", err.Error())
		return
	}
	defer rsp.Body.Close()
	log.Printf("[getPassenger]:  %s\n", string(bodyBytes))
	if rsp.StatusCode == http.StatusOK {
		http_util.CookieChange(conversation, rsp.Cookies())
		pData := bookBean.PresenterVa{}
		err = json.Unmarshal(bodyBytes, &pData)
		if err != nil {
			log.Printf("[getPassenger]: %s\n", err.Error())
			if method == http.MethodPost {
				return GetPassenger(http.MethodPut, conversation, GlobalRepeatSubmitToken)
			} else {
				return
			}
		}
		if pData.Status && pData.Httpstatus == http.StatusOK {
			Passenger = true
			ps = pData.Data.Normal_passengers
		} else {
			log.Printf("[getPassenger]: %v\n", pData)
			Passenger = false
		}
	} else {
		log.Printf("[getPassenger]: %d\n", rsp.StatusCode)
		return
	}
	return
}

func checkOrderInfo(method string, conversation *conversation.Conversation, bookResult *bookBean.BookResult, setType string, u task.TaskDetails) (err error) {
	log.Println("正在检查订单信息...")
	defer func() {
		if re := recover(); re != nil {
			if err == nil {
				err = errors.New(fmt.Sprintf("[checkOrderInfo] %v", re))
				return
			}
		}
	}()
	passengerTicketStr := ""
	oldPassengerStr := ""

	// 拼接passengerTicketStr
	for i, item := range u.TaskPassenger {
		if i == 0 {
			passengerTicketStr += setType + ",0,1," + item.Name + ",1," + item.IdNum + "," + item.TelNum + ",N," + item.AllEncStr
		} else {
			passengerTicketStr += "_" + setType + ",0,1," + item.Name + ",1," + item.IdNum + "," + item.TelNum + ",N," + item.AllEncStr
		}
		oldPassengerStr += item.Name + ",1," + item.IdNum + ",1_"
	}
	//passengerTicketStr := u.SeatTypes + ",0,1," + u. + ",1," + u.Id + "," + u.TelNum + ",N," + "31d2c03567240868c35d68fa9a0d6b5c17cea9706ee43b3a7e066ced20000a692802483a95e936594e91b6096da9c9e8"
	// 拼接oldPassengerStr
	//oldPassengerStr := u.Name + ",1," + u.Id + ",1_"
	// 准备表单数据
	bookResult.PassengerTicketStr = passengerTicketStr
	bookResult.OldPassengerStr = oldPassengerStr
	data := &url.Values{}
	data.Set("cancel_flag", "2")
	data.Set("bed_level_order_num", "000000000000000000000000000000")
	data.Set("passengerTicketStr", passengerTicketStr)
	data.Set("oldPassengerStr", oldPassengerStr)
	data.Set("tour_flag", "dc")
	data.Set("randCode", "")
	data.Set("whatsSelect", "1")
	data.Set("_json_att", "")
	data.Set("REPEAT_SUBMIT_TOKEN", bookResult.GlobalRepeatSubmitToken)
	req, _ := http.NewRequest(method, api.CheckOrderInfo, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	req.Header.Set("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	http_util.SetReqHeader(req)
	http_util.AddReqCookie(conversation.C, req)
	log.Printf("request : %v\n", data)
	rsp, err := conversation.Client.Do(req)
	if err != nil {
		log.Printf("[checkOrderInfo]: %s", err.Error())
		return
	}
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("[checkOrderInfo]: %s", err.Error())
		return
	}
	defer rsp.Body.Close()

	if rsp.StatusCode == http.StatusOK {
		http_util.CookieChange(conversation, rsp.Cookies())
		m := make(map[string]interface{})
		err = json.Unmarshal(bodyBytes, &m)
		if err != nil {
			log.Printf("[checkOrderInfo]: %s", err.Error())
			if method == http.MethodPost {
				return checkOrderInfo(http.MethodPut, conversation, bookResult, setType, u)
			} else {
				return
			}
		}
		//	log.Printf("[checkOrderInfo] bodyBytes : %s\n", string(bodyBytes))
		if m["status"].(bool) && m["data"].(map[string]interface{})["submitStatus"].(bool) {
			bookResult.CheckOrderInfo = true
		} else {
			log.Printf("[checkOrderInfo]: %s\n", m["messages"])
			bookResult.CheckOrderInfo = false
		}
	} else {
		log.Printf("[checkOrderInfo]:  %d\n", rsp.StatusCode)
		return errors.New("network error")
	}
	return
}

/**
 * 该方法为请求请求是否可以加入队列，但响应结果没有用，舍弃请求该页面
 *
 * @param bookResult
 * @return
 */

func getQueueCount(conversation *conversation.Conversation, bookResult *bookBean.BookResult) (err error) {
	log.Println("[getQueueCount] 获取队列信息...")
	defer func() {
		if re := recover(); re != nil {
			if err == nil {
				err = errors.New(fmt.Sprintf("[getQueueCount] %v", re))
				return
			}
		}
	}()
	//Sun Dec 15 2019 00:00:00 GMT+0800 (中国标准时间)
	//  tm:=time.Now()
	data := &url.Values{}
	data.Set("_json_att", "")
	data.Set("from_station_telecode", bookResult.InitDcInfo.FromStationTelecode)
	data.Set("leftTicket", bookResult.InitDcInfo.LeftTicketStr)
	data.Set("purpose_codes", bookResult.InitDcInfo.PurposeCodes)
	data.Set("REPEAT_SUBMIT_TOKEN", bookResult.GlobalRepeatSubmitToken)
	data.Set("seatType", "O")
	data.Set("stationTrainCode", bookResult.InitDcInfo.StationTrainCode)
	data.Set("toStationTelecode", bookResult.InitDcInfo.ToStationTelecode)
	data.Set("train_date", "Sun Dec 15 2019 00:00:00 GMT+0800")
	data.Set("train_location", bookResult.InitDcInfo.TrainLocation)
	data.Set("train_no", bookResult.InitDcInfo.TrainNo)
	req, _ := http.NewRequest(http.MethodPost, api.GetQueueCountURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	req.Header.Set("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	http_util.SetReqHeader(req)
	http_util.AddReqCookie(conversation.C, req)
	log.Printf("request : %v\n", data)
	rsp, err := conversation.Client.Do(req)
	if err != nil {
		log.Printf("[getQueueCount]: %s", err.Error())
		return
	}
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("[getQueueCount]: %s", err.Error())
		return
	}
	defer rsp.Body.Close()
	log.Println(string(bodyBytes))
	if rsp.StatusCode == http.StatusOK {
		http_util.CookieChange(conversation, rsp.Cookies())

	} else {
		log.Printf("[getQueueCount]:  %d\n", rsp.StatusCode)
		return errors.New("network error")
	}
	return
}

/**
 * 该方法为请求进入购票队列
 *
 * @param bookResult
 * @return
 */
func getConfirmSingleForQueue(conversation *conversation.Conversation, bookResult *bookBean.BookResult) (err error) {
	log.Println("[getConfirmSingleForQueue]正在下单...")
	defer func() {
		if re := recover(); re != nil {
			if err == nil {
				err = errors.New(fmt.Sprintf("[getConfirmSingleForQueue] %v", re))
				return
			}
		}
	}()
	data := &url.Values{}
	data.Set("_json_att", "")
	data.Set("choose_seats", "")
	data.Set("dwAll", "N")
	data.Set("key_check_isChange", bookResult.InitDcInfo.KeyCheckIsChange)
	data.Set("leftTicketStr", bookResult.InitDcInfo.LeftTicketStr)
	data.Set("oldPassengerStr", bookResult.OldPassengerStr)
	data.Set("passengerTicketStr", bookResult.PassengerTicketStr)
	data.Set("purpose_codes", bookResult.InitDcInfo.PurposeCodes)
	data.Set("randCode", "")
	data.Set("REPEAT_SUBMIT_TOKEN", bookResult.GlobalRepeatSubmitToken)
	data.Set("roomType", "00")
	data.Set("seatDetailType", "000")
	data.Set("train_location", bookResult.InitDcInfo.TrainLocation)
	data.Set("whatsSelect", "1")
	log.Printf("[getConfirmSingleForQueue]: req : %v\n", data)
	req, _ := http.NewRequest(http.MethodPost, api.ConfirmSingleForQueueURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	//	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	req.Header.Set("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	http_util.SetReqHeader(req)
	http_util.AddReqCookie(conversation.C, req)
	rsp, err := conversation.Client.Do(req)
	if err != nil {
		log.Printf("[getConfirmSingleForQueue]: %s", err.Error())
		return
	}
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("[getConfirmSingleForQueue]: %s", err.Error())
		return
	}
	defer rsp.Body.Close()
	//	log.Printf("[getConfirmSingleForQueue] bodyBytes : %s\n", string(bodyBytes))
	if rsp.StatusCode == http.StatusOK {
		http_util.CookieChange(conversation, rsp.Cookies())
		m := make(map[string]interface{})
		err = json.Unmarshal(bodyBytes, &m)
		if err != nil {
			log.Printf("[getConfirmSingleForQueue]: %s", err.Error())
			return
		}
		if m["status"].(bool) && m["data"].(map[string]interface{})["submitStatus"].(bool) {
			bookResult.ConfirmSingleForQueue = true
		} else {
			log.Printf("[getConfirmSingleForQueue]: %s\n", m["messages"])
			bookResult.CheckOrderInfo = false
		}
	} else {
		log.Printf("[getConfirmSingleForQueue]:  %d\n", rsp.StatusCode)
		return errors.New("network error")
	}
	return

}

/**
 * 该方法为请求排队时间
 *
 * @param bookResult
 * @return
 */
func getQueryOrderTime(conversation *conversation.Conversation, bookResult *bookBean.BookResult) (err error) {
	log.Println("正在查询排队时间...")
	defer func() {
		if re := recover(); re != nil {
			if err == nil {
				err = errors.New(fmt.Sprintf("[getQueryOrderTime] %v", re))
				return
			}
		}
	}()
	data := &url.Values{}
	data.Set("_json_att", "")
	data.Set("random", strconv.FormatInt(time.Now().UnixNano()/10000000, 10))
	data.Set("REPEAT_SUBMIT_TOKEN", bookResult.GlobalRepeatSubmitToken)
	data.Set("tourFlag", "dc")
	req, _ := http.NewRequest(http.MethodPost, api.QueryOrderTime, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	req.Header.Set("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	http_util.SetReqHeader(req)
	http_util.AddReqCookie(conversation.C, req)

	//等待三分钟
	i := 0
	Q := bookBean.QueryTimeResult{}
	for (!Q.OK) || (i > 300) {
		Q = getQueryOrderTimeMethod(req, conversation, bookResult)
		time.Sleep(time.Second * 5)
		i += 5
	}
	if Q.OK {
		log.Println("订票成功,订单号：" + Q.OrderId)
	}
	bookResult.QueryOrderTime = Q.OK
	bookResult.QueryTimeResult = Q
	return
}

/**
 * 该方法为请求排队时间的方法体，因为要循环请求，所以单独拿出来作为一个方法
 *
 * @param getQueryOrderTime
 * @param bookResult
 * @return
 */
func getQueryOrderTimeMethod(r *http.Request, conversation *conversation.Conversation, bookResult *bookBean.BookResult) (queryTimeResult bookBean.QueryTimeResult) {
	queryTimeResult = bookBean.QueryTimeResult{}
	defer func() {
		if re := recover(); re != nil {
			log.Println("[getQueryOrderTimeMethod] recover=" + fmt.Sprintf("%s", re))
		}
	}()

	rsp, err := conversation.Client.Do(r)
	if err != nil {
		log.Printf("[getQueryOrderTimeMethod]: %s\n", err.Error())
		return
	}
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("[getQueryOrderTimeMethod]: %s\n", err.Error())
		return
	}
	defer rsp.Body.Close()
	log.Printf("[getQueryOrderTimeMethod] bodyBytes : %s\n", string(bodyBytes))
	if rsp.StatusCode == http.StatusOK {
		http_util.CookieChange(conversation, rsp.Cookies())
		m := make(map[string]interface{})
		err = json.Unmarshal(bodyBytes, &m)
		if err != nil {
			log.Printf("[getQueryOrderTimeMethod]: %s\n", err.Error())
			return
		}
		// 出现错误时提示错误信息
		if m["data"].(map[string]interface{})["errorcode"] != nil {
			log.Printf("[getQueryOrderTimeMethod]: %s\n", m["data"].(map[string]interface{})["msg"])
			return
		}
		if m["status"].(bool) && m["data"].(map[string]interface{})["queryOrderWaitTimeStatus"].(bool) {
			queryTimeResult.WaitTime = m["data"].(map[string]interface{})["waitTime"].(float64)
			queryTimeResult.WaitCount = m["data"].(map[string]interface{})["waitCount"].(float64)
			if m["data"].(map[string]interface{})["orderId"] != nil {
				queryTimeResult.OK = true
				queryTimeResult.OrderId = m["data"].(map[string]interface{})["orderId"].(string)
			}
			log.Printf("[getQueryOrderTimeMethod]: 等待时间:%f,等待人数:%f\n", queryTimeResult.WaitTime, queryTimeResult.WaitCount)
		} else {
			log.Printf("[getQueryOrderTimeMethod]: %s\n", m["messages"])
			bookResult.CheckOrderInfo = false
		}
	} else {
		log.Printf("[getQueryOrderTimeMethod]:  %d\n", rsp.StatusCode)
		return
	}
	return
}

/**
 * 请求resultOrderForQueue页面，虽然该页面不会返回任何信息，但是如果不请求该页面，不能请求后面的订单信息页面
 *
 * @param bookResult
 * @return
 */
func getResultOrderForQueue(conversation *conversation.Conversation, bookResult *bookBean.BookResult) (err error) {
	log.Println("请求resultOrderForQueue页面...")
	defer func() {
		if re := recover(); re != nil {
			log.Println("[getResultOrderForQueue] recover=" + fmt.Sprintf("%v", re))
		}
	}()

	data := &url.Values{}
	data.Set("_json_att", "")
	data.Set("orderSequence_no", bookResult.QueryTimeResult.OrderId)
	data.Set("REPEAT_SUBMIT_TOKEN", bookResult.GlobalRepeatSubmitToken)
	req, _ := http.NewRequest(http.MethodPost, api.ResultOrderForQueue, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	req.Header.Set("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	http_util.SetReqHeader(req)
	http_util.AddReqCookie(conversation.C, req)
	rsp, err := conversation.Client.Do(req)
	if err != nil {
		log.Printf("[getResultOrderForQueue]: %s", err.Error())
		return
	}
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("[getResultOrderForQueue]: %s", err.Error())
		return
	}
	defer rsp.Body.Close()
	//	log.Printf("[getResultOrderForQueue] bodyBytes : %s\n", string(bodyBytes))
	if rsp.StatusCode == http.StatusOK {
		http_util.CookieChange(conversation, rsp.Cookies())
		m := make(map[string]interface{})
		err = json.Unmarshal(bodyBytes, &m)
		if err != nil {
			log.Printf("[getResultOrderForQueue]: %s", err.Error())
			return
		}
		if m["status"].(bool) && m["data"].(map[string]interface{})["submitStatus"].(bool) {
			bookResult.ResultOrderForQueue = true
		} else {
			log.Printf("[checkOrderInfo]: %s\n", m["messages"])
			bookResult.ResultOrderForQueue = false
		}
	} else {
		log.Printf("[getResultOrderForQueue]:  %d\n", rsp.StatusCode)
		return errors.New("network error")
	}
	return
}

/**
 * 获取订单信息
 * 该方法中的请求地址URL需要拼接上当前时间的时间戳
 *
 * @param bookResult
 * @return
 */
func getOrderMsg(conversation *conversation.Conversation, bookResult *bookBean.BookResult) (err error) {
	log.Println("请求getOrderMsg页面...")
	defer func() {
		if re := recover(); re != nil {
			log.Println("[getOrderMsg] recover=" + fmt.Sprintf("%v", re))
		}
	}()
	data := &url.Values{}
	data.Set("_json_att", "")
	data.Set("REPEAT_SUBMIT_TOKEN", bookResult.GlobalRepeatSubmitToken)
	req, _ := http.NewRequest(http.MethodPost, api.BookResult, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	req.Header.Set("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	http_util.SetReqHeader(req)
	http_util.AddReqCookie(conversation.C, req)
	rsp, err := conversation.Client.Do(req)
	if err != nil {
		log.Printf("[getOrderMsg]: %s", err.Error())
		return
	}
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("[getOrderMsg]: %s", err.Error())
		return
	}
	defer rsp.Body.Close()
	//	log.Printf("[getResultOrderForQueue] bodyBytes : %s\n", string(bodyBytes))
	if rsp.StatusCode == http.StatusOK {
		http_util.CookieChange(conversation, rsp.Cookies())
		str := string(bodyBytes)
		htmls := strings.Split(str, "\n")
		var result = ""
		for _, line := range htmls {
			if strings.Contains(line, "var passangerTicketList") {
				result = line[strings.Index(line, "[") : len(line)-1]
			}
		}
		if len(result) > 0 {
			log.Println("订单已完成，请登录12306查看")
			bookResult.Finish = true
			//	bookResult.setOrderMsg(msg);

		} else {
			log.Println("订单信息查询失败，可能该订单已完成，请登录12306查看，")
		}

	} else {
		log.Printf("[getOrderMsg]:  %s\n", string(bodyBytes))
		return errors.New("network error")
	}
	return
}
