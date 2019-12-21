package service

import (
	"book-query-srv/stations"
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/basic/utils/conversation"
	"gitee.com/qianxunke/book-ticket-common/basic/utils/http_util"
	"gitee.com/qianxunke/book-ticket-common/proto/ticket"
	"gitee.com/qianxunke/book-ticket-common/ticket/static/api"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

/**
 * 查询符合出发站、终点站、出发日期的列车信息的方法
 * @param u
 *   用户信息
 * @return
 */
var (
	conversation2 *conversation.Conversation
)

func InitCon() {
	conversation2 = &conversation.Conversation{}
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
}

func (s *service) queryTrainMessage(qq string, que ticket.In_GetTrainInfoList) (tran []*ticket.Train, err error) {
	if conversation2 == nil {
		InitCon()
	}
	//ADULT
	req1, _ := http.NewRequest(http.MethodGet, api.Query+"Z?leftTicketDTO.train_date="+que.TrainDate+"&leftTicketDTO.from_station="+que.FindFrom+"&leftTicketDTO.to_station="+que.FindTo+"&purpose_codes="+que.PurposeCodes, nil)
	http_util.AddReqCookie(conversation2.C, req1)
	http_util.SetReqHeader(req1)
	rsp1, err := conversation2.Client.Do(req1)
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
	if rsp1.StatusCode == http.StatusOK {
		queryItem := &QueryItem{}
		err = json.Unmarshal(str, &queryItem)
		if err != nil {
			log.Printf("[QueryTrainMessage] error %v", string(str))
			return tran, err
		}
		if len(queryItem.Data.Result) > 0 {
			trans, err := formatQueryMessage(queryItem.Data.Result)
			if err != nil {
				log.Printf("[QueryTrainMessage] error %v", err)
				return tran, err
			}
			return trans, nil
		}
	} else {
		//log.Printf("[QueryTrainMessage] net error %d :, %s", rsp.StatusCode, string(str))
		err = errors.New("[QueryTrainMessage] net error")
		return tran, err
	}
	err = errors.New("[QueryTrainMessage] net error")
	return tran, err

}

/*
function b4(ct, cv) {
        var cs = [];
        for (var cr = 0; cr < ct.length; cr++) {
            var cw = [];
            var cq = ct[cr].split("|"); //用"|"进行字符串切割
            cw.secretHBStr = cq[36];
            cw.secretStr = cq[0]; //secretStr索引为0
            cw.buttonTextInfo = cq[1];
            var cu = [];
            cu.train_no = cq[2]; //车票号
            cu.station_train_code = cq[3]; //车次
            cu.start_station_telecode = cq[4]; //起始站代号
            cu.end_station_telecode = cq[5]; //终点站代号
            cu.from_station_telecode = cq[6]; //出发站代号
            cu.to_station_telecode = cq[7]; //到达站代号
            cu.start_time = cq[8]; //出发时间
            cu.arrive_time = cq[9]; //到达时间
            cu.lishi = cq[10]; //历时
            cu.canWebBuy = cq[11]; //是否能购买：Y 可以
            cu.yp_info = cq[12];
            cu.start_train_date = cq[13]; //出发日期
            cu.train_seat_feature = cq[14];
            cu.location_code = cq[15];
            cu.from_station_no = cq[16];
            cu.to_station_no = cq[17];
            cu.is_support_card = cq[18];
            cu.controlled_train_flag = cq[19];
            cu.gg_num = cq[20] ? cq[20] : "--";
            cu.gr_num = cq[21] ? cq[21] : "--";
            cu.qt_num = cq[22] ? cq[22] : "--";
            cu.rw_num = cq[23] ? cq[23] : "--"; //软卧
            cu.rz_num = cq[24] ? cq[24] : "--"; //软座
            cu.tz_num = cq[25] ? cq[25] : "--";
            cu.wz_num = cq[26] ? cq[26] : "--"; //无座
            cu.yb_num = cq[27] ? cq[27] : "--";
            cu.yw_num = cq[28] ? cq[28] : "--"; //硬卧
            cu.yz_num = cq[29] ? cq[29] : "--";
            cu.ze_num = cq[30] ? cq[30] : "--"; //二等座
            cu.zy_num = cq[31] ? cq[31] : "--"; //一等座
            cu.swz_num = cq[32] ? cq[32] : "--"; //商务特等座
            cu.srrb_num = cq[33] ? cq[33] : "--";
            cu.yp_ex = cq[34];
            cu.seat_types = cq[35];
            cu.exchange_train_flag = cq[36];
            cu.from_station_name = cv[cq[6]];
            cu.to_station_name = cv[cq[7]];
            cw.queryLeftNewDTO = cu;
            cs.push(cw)
        }
        return cs
    }
商务特等座：32
一等座：31
二等座：30
高级软卧：21
软卧：23
动卧：33
硬卧：28
软座：24
硬座：29
无座：26
其他：22
备注：1

start_train_date:车票出发日期：13
*/
/**
 * 格式化输出列车信息
 *
 * @param s
 * @param stationCode
 * @return
 */
func formatQueryMessage(s []string) (trans []*ticket.Train, err error) {
	defer func() {
		if re := recover(); re != nil {
			err = errors.New(fmt.Sprintf("%v", re))
		}
	}()
	if len(s) <= 0 {
		return
	}
	for i := 0; i < len(s); i++ {
		ss := strings.Split(s[i], "|")
		t := ticket.Train{}
		t.SecretStr = ss[0]
		t.Bz = ss[1]
		t.TrainCode = ss[2]
		t.Num = ss[3]
		t.From = ss[4]
		t.To = ss[5]
		t.FindFrom = ss[6]
		t.FindTo = ss[7]
		t.StartTime = ss[8]
		t.EndTime = ss[9]
		t.CostTime = ss[10]
		t.CanBuy = ss[11]
		t.TrainDate = ss[13]
		t.Wz = ss[26]
		t.Yz = ss[29]
		t.Rz = ss[24]
		t.Yw = ss[28]
		t.Dw = ss[33]
		t.Rw = ss[23]
		t.Gjrw = ss[21]
		t.Edz = ss[30]
		t.Ydz = ss[31]
		t.Swtdz = ss[32]
		trans = append(trans, &t)
	}
	return
}

/**
 * 输出列车信息
 * @param newMessage
 * @param stationCode
 */
func printTrainMessage(trans []*ticket.Train) {
	fmt.Print("序号\t")
	fmt.Print("车次\t")
	fmt.Print("始发站\t")
	fmt.Print("终到站\t")
	fmt.Print("查询始发站\t")
	fmt.Print("查询终点站\t")
	fmt.Print("出发时间\t")
	fmt.Print("到站时间\t")
	fmt.Print("历时\t")
	fmt.Print("商务特等座\t")
	fmt.Print("一等座\t")
	fmt.Print("二等座\t")
	fmt.Print("高级软卧\t")
	fmt.Print("软卧\t")
	fmt.Print("动卧\t")
	fmt.Print("硬卧\t")
	fmt.Print("软座\t")
	fmt.Print("硬座\t")
	fmt.Print("无座\t")
	fmt.Print("备注\t")
	fmt.Print("出发日期\t")
	fmt.Println("可否购票\t")

	for i := 0; i < len(trans); i++ {
		fmt.Printf("%d\t", i)
		fmt.Print(trans[i].Num + "\t")
		fmt.Print(stations.GetStationValueByKey(trans[i].From) + "\t")
		fmt.Print(stations.GetStationValueByKey(trans[i].To) + "\t")
		fmt.Print(stations.GetStationValueByKey(trans[i].FindFrom) + "\t")
		fmt.Print(stations.GetStationValueByKey(trans[i].FindTo) + "\t")
		fmt.Print(trans[i].StartTime + "\t")
		fmt.Print(trans[i].EndTime + "\t")
		fmt.Print(trans[i].CostTime + "\t")
		fmt.Print(trans[i].Swtdz + "\t")
		fmt.Print(trans[i].Ydz + "\t")
		fmt.Print(trans[i].Edz + "\t")
		fmt.Print(trans[i].Gjrw + "\t")
		fmt.Print(trans[i].Rw + "\t")
		fmt.Print(trans[i].Dw + "\t")
		fmt.Print(trans[i].Yw + "\t")
		fmt.Print(trans[i].Rz + "\t")
		fmt.Print(trans[i].Yz + "\t")
		fmt.Print(trans[i].Wz + "\t")
		fmt.Print(trans[i].Bz + "\t")
		fmt.Print(trans[i].TrainDate + "\t")
		fmt.Print(trans[i].CanBuy + "\t")
		fmt.Println()
	}
}

type QueryItem struct {
	Data       Data
	Httpstatus int
	Messages   string
	status     bool
}

type Data struct {
	Flag   string
	Map    map[string]string
	Result []string
}
