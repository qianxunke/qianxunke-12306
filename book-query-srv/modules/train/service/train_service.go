package service

import (
	"fmt"
	ticketProto "gitee.com/qianxunke/book-ticket-common/proto/ticket"
	"sync"
)

type service struct {
}

var (
	s *service
	m sync.RWMutex
)

type Service interface {
	//获取信息
	GetTrainById(req *ticketProto.In_GetTrainInfo) (rsp *ticketProto.Out_GetTrainInfo)
	//修改信息
	UpdateTrainInfo(req *ticketProto.In_UpdateTrainInfo) (rsp *ticketProto.Out_UpdateTrainInfo)
	//获取列表
	GetTrainInfoList(req *ticketProto.In_GetTrainInfoList) (rsp *ticketProto.Out_GetTrainInfoList)
	//新建信息
	CreateTrain(req *ticketProto.In_AddTrain) (rsp *ticketProto.Out_AddTrain)
}

//获取服务
func GetService() (Service, error) {
	if s == nil {
		return nil, fmt.Errorf("[GetService] GetService 未初始化")
	}
	return s, nil
}

func Init() {
	m.Lock()
	defer m.Unlock()
	if s != nil {
		return
	}
	s = &service{}
}
