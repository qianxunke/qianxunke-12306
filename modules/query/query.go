package query

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"qianxunke-12306/common/http_util"
	"qianxunke-12306/config/api"
	"qianxunke-12306/modules/conversation"
	"qianxunke-12306/modules/login"
	"qianxunke-12306/modules/query/bean"
	"qianxunke-12306/modules/sysinit/stations"
)

/**
 * 查询符合出发站、终点站、出发日期的列车信息的方法
 * @param u
 *   用户信息
 * @return
 */
func QueryTrainMessage(uN string, conversation *conversation.Conversation) {
	u := &login.User{RideDate: "2019-12-13", Departure: "北京", Terminus: "上海"}
	req, _ := http.NewRequest(http.MethodGet, api.Query+"?leftTicketDTO.train_date="+u.RideDate+"&leftTicketDTO.from_station="+stations.GetStationValueByKey(u.Departure)+"&leftTicketDTO.to_station="+stations.GetStationValueByKey(u.Terminus)+"&purpose_codes=ADULT", nil)
	http_util.AddReqCookie(conversation.C, req)
	http_util.SetReqHeader(req)
	rsp, err := conversation.Client.Do(req)
	if err != nil {
		log.Printf("[QueryTrainMessage] error %v", err)
		return
	}
	str, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("[QueryTrainMessage] error %v", err)
		return
	}
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusOK {
		queryItem := &bean.QueryItem{}
		err = json.Unmarshal(str, &queryItem)
		if err != nil {
			log.Printf("[QueryTrainMessage] error %v", err)
			return
		}

		log.Printf("result %+v", uN)

	} else {
		log.Printf("[QueryTrainMessage] net error %d :, %v", rsp.StatusCode, err)
	}
	return

}
