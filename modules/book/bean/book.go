package bean

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
)

type BookResult struct {
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
	log.Printf("[formatInewTicketInfoForPassengerFormnitDc] doc :%s\n", newTicketInfoForPassengerForm)
	log.Printf("[formatInitDc] doc :%v\n", doc)
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
