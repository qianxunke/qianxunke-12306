package book_service

import (
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/proto/task"
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
	GetTaskById(req *task.In_GetTaskInfo) (rsp *task.Out_GetTaskInfo)
	//修改信息
	UpdateBrandInfo(req *task.In_UpdateTaskInfo) (rsp *task.Out_UpdateTaskInfo)
	//获取列表
	GetBrands(req *task.Out_GetTaskInfoList) (rsp *task.Out_GetTaskInfoList)

	//新建信息
	CreateTask(req *task.In_AddTask) (rsp *task.Out_AddTask)
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
