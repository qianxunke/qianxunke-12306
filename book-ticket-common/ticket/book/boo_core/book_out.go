package boo_core

import (
	"gitee.com/qianxunke/book-ticket-common/basic/utils/conversation"
	"gitee.com/qianxunke/book-ticket-common/proto/task"
	bookBean "gitee.com/qianxunke/book-ticket-common/ticket/book/bean"
	"gitee.com/qianxunke/book-ticket-common/ticket/login"
	"gitee.com/qianxunke/book-ticket-common/ticket/query/bean"
	"log"
	"net/http"
)

//下单
func Book(conversation conversation.Conversation, chooseTrain bean.Train, u task.TaskDetails) (ok bool) {
	// 设置绑定会话的客户端
	boolResult := &bookBean.BookResult{}
	boolResult.SelectTran = chooseTrain
	boolResult.Task = u
	// 检查登录状态
	err := login.CheckUserStatus(boolResult, &conversation)
	if err != nil {
		log.Println("[login.CheckUserStatus] error :" + err.Error())
		return
	}
	if !boolResult.CheckUser {
		log.Println("用户未登录...")
		return
	}
	err = submitOrder(&conversation, boolResult)
	if err != nil {
		log.Println("[submitOrder] error :" + err.Error())
		return
	}
	if !boolResult.SubmitOrder {
		log.Println("有未完成订单...")
		return
	}
	// 从InitDc获取必要信息
	boolResult.GlobalRepeatSubmitToken, boolResult.InitDcInfo, err = GetInitDc(&conversation)
	if err != nil {
		log.Println("[getInitDc] error :" + err.Error())
		return
	}
	if len(boolResult.GlobalRepeatSubmitToken) <= 0 {
		log.Println("InitDC请求失败...")
		return
	}
	// 获取乘客信息
	boolResult.Passenger, _, err = GetPassenger(http.MethodPost, &conversation, boolResult.GlobalRepeatSubmitToken)
	if err != nil {
		log.Println("[getPassenger] error :" + err.Error())
		return
	}
	if !boolResult.Passenger {
		log.Println("乘客信息请求失败...")
		return
	}
	err = checkOrderInfo(http.MethodPost, &conversation, boolResult, u)
	if err != nil {
		log.Println("[checkOrderInfo] error :" + err.Error())
		return
	}
	if !boolResult.CheckOrderInfo {
		// bookResult=BookUtils.getQueueCount(bookResult);
		log.Println("订单信息检查错误...")
		return
	}
	_ = getQueueCount(&conversation, boolResult)

	err = getConfirmSingleForQueue(&conversation, boolResult)
	if err != nil {
		log.Println("[getConfirmSingleForQueue] error :" + err.Error())
		return
	}
	if !boolResult.ConfirmSingleForQueue {
		log.Println("下单失败...")
		return
	}
	err = getQueryOrderTime(&conversation, boolResult)
	if err != nil {
		log.Println("[getQueryOrderTime] error :" + err.Error())
		return
	}
	if !boolResult.QueryOrderTime {
		log.Println("查询排队时间失败...")
		return
	}
	err = getResultOrderForQueue(&conversation, boolResult)
	if err != nil {
		log.Println("[getResultOrderForQueue] error :" + err.Error())
		return
	}
	err = getOrderMsg(&conversation, boolResult)
	/*
		if !boolResult.QueryOrderTime {
			log.Println("查询排队时间失败...")
			return
		}

	*/
	if boolResult.Finish {
		log.Println("恭喜您，订票成功，请在30分钟内登录12306完成支付！")
	}
	ok = true
	return
}
