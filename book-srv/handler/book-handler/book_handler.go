package book_handler

import (
	"book-srv/modules/book/book_service"
	"context"
	"gitee.com/qianxunke/book-ticket-common/proto/task"
	"log"
)

type Handler struct{}

var (
	s book_service.Service
)

func Init() {
	_userService, err := book_service.GetService()
	if err != nil {
		log.Fatal("[Init] 初始化Handler错误")
		return
	}
	s = _userService
}

func (e *Handler) GetTaskInfo(ctx context.Context, req *task.In_GetTaskInfo, rsp *task.Out_GetTaskInfo) error {
	log.Println("Received UserInfo.GetUserInfo request")
	response := s.GetTaskById(req)
	rsp.Error = response.Error
	rsp.TaskDetails = response.TaskDetails
	return nil
}

func (e *Handler) GetTaskInfoList(ctx context.Context, req *task.In_GetTaskInfoList, rsp *task.Out_GetTaskInfoList) error {
	log.Println("Received UserInfo.GetUserInfoList request")
	response := s.GetTasks(req)
	rsp.Error = response.Error
	rsp.TaskDetailsList = response.TaskDetailsList
	rsp.Pages = response.Pages
	rsp.Limit = response.Limit
	rsp.Total = response.Total
	return nil
}

func (e *Handler) UpdateTaskInfo(ctx context.Context, req *task.In_UpdateTaskInfo, rsp *task.Out_UpdateTaskInfo) error {
	log.Printf("Received Task.UpdateTaskInfo request : %v", req)
	response := s.UpdateTaskInfo(req)
	rsp.Error = response.Error
	return nil
}

func (e *Handler) AddTask(ctx context.Context, req *task.In_AddTask, rsp *task.Out_AddTask) error {
	log.Println("Received Task.AddTask request")
	response := s.CreateTask(req)
	rsp.Error = response.Error
	rsp.TaskDetails = response.TaskDetails
	return nil
}

func (e *Handler) GetUserTaskList(ctx context.Context, req *task.In_GetUserTaskList, rsp *task.Out_GetTaskInfoList) error {
	log.Println("Received Task.GetUserTaskList request")
	response := s.GetUserTaskList(req)
	rsp.Error = response.Error
	rsp.Total = response.Total
	rsp.TaskDetailsList = response.TaskDetailsList
	return nil
}

func (e *Handler) UpdateTaskStatus(ctx context.Context, req *task.In_UpdateTaskStatus, rsp *task.Out_UpdateTaskStatus) error {
	log.Println("Received Task.UpdateTaskStatus request")
	response := s.UpdateTaskStatus(req)
	rsp.Error = response.Error
	return nil
}
