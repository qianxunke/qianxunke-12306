package train

import (
	"book-query-srv/modules/train/service"
	"context"
	"gitee.com/qianxunke/book-ticket-common/proto/ticket"
	"log"
)

type Handler struct{}

var (
	brandS service.Service
)

func Init() {
	_brandS, err := service.GetService()
	if err != nil {
		log.Fatal("[Init] 初始化BrandHandler错误: " + err.Error())
		return
	}
	brandS = _brandS
}

func (e *Handler) GetTrainInfo(ctx context.Context, req *ticket.In_GetTrainInfo, rsp *ticket.Out_GetTrainInfo) error {
	log.Println("Received GetTrainInfo request")
	response := brandS.GetTrainById(req)
	rsp.Error = response.Error
	rsp.Train = response.Train
	return nil
}
func (e *Handler) UpdateTrainInfo(ctx context.Context, req *ticket.In_UpdateTrainInfo, rsp *ticket.Out_UpdateTrainInfo) error {
	log.Println("Received UpdateTrainInfo request")
	response := brandS.UpdateTrainInfo(req)
	rsp.Error = response.Error
	return nil
}
func (e *Handler) GetTrainInfoList(ctx context.Context, req *ticket.In_GetTrainInfoList, rsp *ticket.Out_GetTrainInfoList) error {
	log.Println("Received GetTrainInfoList request")
	response := brandS.GetTrainInfoList(req)
	rsp.TrainList = response.TrainList
	rsp.Pages = response.Pages
	rsp.Limit = response.Limit
	rsp.Error = response.Error
	rsp.Total = response.Total
	return nil
}

func (e *Handler) AddTrain(ctx context.Context, req *ticket.In_AddTrain, rsp *ticket.Out_AddTrain) error {
	log.Println("Received AddTrain request")
	response := brandS.CreateTrain(req)
	rsp.Error = response.Error
	return nil
}
