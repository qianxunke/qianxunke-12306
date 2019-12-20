package book_service

import (
	"book-srv/modules/book/book_dao"
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/proto/task"
	"log"
	"net/http"
	"strings"
)

//获取信息
func (s *service) GetTaskById(req *task.In_GetTaskInfo) (rsp *task.Out_GetTaskInfo) {
	rsp = &task.Out_GetTaskInfo{}
	rsp.Error = &task.Error{}
	if len(req.TaskId) <= 0 {
		rsp.Error = &task.Error{
			Code:    http.StatusBadRequest,
			Message: "请求参数有误！",
		}
		return
	}
	dao, err := book_dao.GetDao()
	if err != nil {
		rsp.Error = &task.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	rsp.TaskDetails, err = dao.FindById(req.TaskId)
	if err != nil {
		rsp.Error = &task.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	rsp.Error = &task.Error{
		Code:    http.StatusOK,
		Message: "查询成功！",
	}
	return

}

//修改信息
func (s *service) UpdateTaskInfo(req *task.In_UpdateTaskInfo) (rsp *task.Out_UpdateTaskInfo) {
	defer func() {
		if re := recover(); re != nil {
			rsp.Error = &task.Error{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("[UpdateBrandInfo] error %v", re),
			}
		}
	}()

	rsp = &task.Out_UpdateTaskInfo{}
	dao, err := book_dao.GetDao()
	err = dao.Update(req.TaskDetails)
	if err != nil {
		rsp.Error = &task.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	rsp.Error = &task.Error{
		Code:    http.StatusOK,
		Message: "修改成功！",
	}
	return
}

//获取列表
func (s *service) GetTasks(req *task.In_GetTaskInfoList) (rsp *task.Out_GetTaskInfoList) {
	rsp = &task.Out_GetTaskInfoList{}
	//对参数鉴权
	if req.Limit == 0 {
		req.Limit = 10 //默认10个分页
	}
	if req.Limit > 1000 { //每一页数量
		req.Limit = 1000
	}
	if req.Pages <= 0 { //页数
		req.Pages = 1
	}
	var err error
	orderByStr := "created_time DESC"
	dao, err := book_dao.GetDao()
	if err != nil {
		log.Printf("ERROR: %v", err)
		rsp.Error = &task.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	rsp, err = dao.SimpleQuery(req.Limit, req.Pages, 1, req.SearchKey, req.StartTime, req.EndTime, orderByStr)
	if err != nil {
		log.Printf("ERROR: %v", err)
		rsp.Error = &task.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	var message string
	if rsp.Total > 0 && len(rsp.TaskDetailsList) > 0 {
		message = "查询成功！"
	} else {
		message = "没有数据了！"
	}
	//统计有多少条
	rsp.Error = &task.Error{
		Code:    http.StatusOK,
		Message: message,
	}
	rsp.Limit = req.Limit
	rsp.Pages = req.Pages
	return

}

//新建信息
func (s *service) CreateTask(req *task.In_AddTask) (rsp *task.Out_AddTask) {
	rsp = &task.Out_AddTask{}
	//查询该等级是否存在
	dao, err := book_dao.GetDao()
	if err != nil {
		log.Printf("ERROR: %v", err)
		rsp.Error = &task.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	//判断当前任务是否可以直接开始
	isCan := false
	d_temp := strings.Split(req.TaskDetails.Task.TrainDates, ",")
	for _, item := range d_temp {
		d := isCanQuery(item)
		if d > 0 && d < 30 {
			isCan = true
		}
	}
	if isCan {
		req.TaskDetails.Task.Status = 1
	} else {
		req.TaskDetails.Task.Status = 4
	}
	err = dao.Insert(req.TaskDetails)
	if err != nil {
		rsp.Error = &task.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	rsp.Error = &task.Error{
		Code:    http.StatusOK,
		Message: "新增成功!",
	}
	return
}
func (s *service) GetNeedTicketList(limit int64, pages int64, status int64) (rsp []task.Task, err error) {

	//对参数鉴权
	if limit == 0 {
		limit = 10 //默认10个分页
	}
	if limit > 1000 { //每一页数量
		limit = 1000
	}
	if pages <= 0 { //页数
		pages = 1
	}
	dao, err := book_dao.GetDao()
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}
	rsp, err = dao.TicketQuery(limit, pages, status)
	return

}

//获取用户列表
func (s *service) GetUserTaskList(req *task.In_GetUserTaskList) (rsp *task.Out_GetTaskInfoList) {
	dao, err := book_dao.GetDao()
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}
	rsp, err = dao.GetUserTask(req.UserId)
	return
}

//修改信息
func (s *service) UpdateTaskStatus(req *task.In_UpdateTaskStatus) (rsp *task.Out_UpdateTaskStatus) {

	dao, err := book_dao.GetDao()
	rsp = &task.Out_UpdateTaskStatus{}
	if err != nil {
		log.Printf("ERROR: %v", err)
		rsp.Error = &task.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		return
	}
	err = dao.UpdateStatus(req.TaskId, req.Status)
	if err != nil {
		rsp.Error = &task.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	} else {
		rsp.Error = &task.Error{Code: http.StatusOK, Message: "更新成功"}
	}
	return
}
