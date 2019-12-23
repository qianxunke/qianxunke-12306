package bean

import (
	"encoding/json"
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/proto/task"
	"gitee.com/qianxunke/book-ticket-common/ticket/query/bean"
	"log"
	"net/url"
	"strings"
)

type BookResult struct {
	Task task.TaskDetails
	//选择的车次
	SelectTran              bean.Train
	CheckUser               bool
	SubmitOrder             bool
	GlobalRepeatSubmitToken string
	InitDcInfo              InitDC
	InitDc                  bool
	Passenger               bool
	CheckOrderInfo          bool
	PassengerTicketStr      string
	OldPassengerStr         string
	ConfirmSingleForQueue   bool
	QueryOrderTime          bool
	QueryTimeResult         QueryTimeResult
	ResultOrderForQueue     bool
	FormatBookResult        bool
	OrderMsg                OrderMsg
	Finish                  bool
}

/**
 * 因为InitDc中信息较多，单独开发一个类
 *
 */
type InitDC struct {
	TrainDateTime       string
	FromStationTelecode string
	LeftTicketStr       string
	PurposeCodes        string
	StationTrainCode    string
	ToStationTelecode   string
	TrainLocation       string
	TrainNo             string
	LeftDetails         []string
	KeyCheckIsChange    string
}

type QueryTimeResult struct {
	OK        bool
	WaitTime  float64
	WaitCount float64
	OrderId   string
}

/**
 * 该类为订票成功后生成的订单信息类
 */
type OrderMsg struct {
	SequenceNo          string
	PassengerIdTypeName string
	PassengerName       string
	PassengerIdNo       string
	FromStationName     string
	ToStationName       string
	StationTrainCode    string
	StartTrainDate      string
	TicketPrice         string
	TicketTypeName      string
	CoachName           string
	SeatName            string
	SeatTypeName        string
}

type PresenterVa struct {
	Status     bool
	Httpstatus int
	Data       struct {
		Notify_for_gat    string
		IsExist           bool
		ExMsg             string
		Two_isOpenClick   []string
		Other_isOpenClick []string
		Dj_passengers     interface{}
		Normal_passengers []Normal_passengers
	}
	messages         interface{}
	validateMessages interface{}
}

/**
 * 格式化InitDc的方法
 *
 * @param ticketInfoForPassengerForm
 *            返回的html源码
 * @return
 */
func FormatInitDc(ticketInfoForPassengerForm string) (intDoc InitDC, err error) {
	intDoc = InitDC{}
	newTicketInfoForPassengerForm := strings.ReplaceAll(ticketInfoForPassengerForm, "'", "\"")
	type Doc struct {
		OrderRequestDTO struct {
			Train_date struct {
				Time int64
			}
			From_station_telecode string
			Station_train_code    string
			To_station_telecode   string
		}
		QueryLeftTicketRequestDTO struct {
			Train_no string
		}
		LeftTicketStr      string
		Purpose_codes      string
		Train_location     string
		LeftDetails        []string
		Key_check_isChange string `json:"key_check_isChange"`
	}
	doc := &Doc{}
	err = json.Unmarshal([]byte(newTicketInfoForPassengerForm), &doc)
	if err != nil {
		return
	}
	//	log.Printf("[formatInewTicketInfoForPassengerFormnitDc] doc :%s\n", newTicketInfoForPassengerForm)
	//	log.Printf("[formatInitDc] doc :%v\n", doc)
	intDoc.TrainDateTime = fmt.Sprintf("%d", doc.OrderRequestDTO.Train_date.Time)
	intDoc.FromStationTelecode = doc.OrderRequestDTO.From_station_telecode
	t, _ := url.PathUnescape(doc.LeftTicketStr)
	intDoc.LeftTicketStr = t
	intDoc.PurposeCodes = doc.Purpose_codes
	intDoc.StationTrainCode = doc.OrderRequestDTO.Station_train_code
	intDoc.ToStationTelecode = doc.OrderRequestDTO.To_station_telecode
	intDoc.TrainLocation = doc.Train_location
	intDoc.TrainNo = doc.QueryLeftTicketRequestDTO.Train_no
	intDoc.LeftDetails = doc.LeftDetails
	intDoc.KeyCheckIsChange = doc.Key_check_isChange
	log.Printf("[formatInitDc] intDoc :%v\n", intDoc)
	return
}

type Normal_passengers struct {
	Passenger_name         string //:"王芳平",
	Sex_code               string //":"M",
	Sex_name               string //":"男",
	Born_date              string //:"1980-01-01 00:00:00",
	Country_code           string //:"CN",
	Passenger_id_type_code string //:"1",
	Passenger_id_type_name string //:"中国居民身份证",
	Passenger_id_no        string //":"5225***********610",
	Passenger_type         string //:"1",
	Passenger_flag         string //:"0",
	Passenger_type_name    string //:"成人",
	Mobile_no              string //:"18334142052",
	Phone_no               string //":"",
	Email                  string //:"",
	Address                string //:"",
	Postalcode             string //":"",
	First_letter           string //:"WFP",
	RecordCount            string //":"4",
	Total_times            string //":"99",
	Index_id               string //:"1",
	AllEncStr              string //":"31d2c03567240868c35d68fa9a0d6b5c17cea9706ee43b3a7e066ced20000a692802483a95e936594e91b6096da9c9e8",
	IsAdult                string //":"Y",
	IsYongThan10           string //":"N",
	IsYongThan14           string //":"N",
	IsOldThan60            string //   :"N",
	Gat_born_date          string //":"",
	Gat_valid_date_start   string //''":"",
	Gat_valid_date_end     string //:"",
	Gat_version            string //:""
}
