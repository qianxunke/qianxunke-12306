package task_handler

import (
	"context"
	"gitee.com/qianxunke/book-ticket-common/basic"
	"gitee.com/qianxunke/book-ticket-common/basic/api_common"
	"gitee.com/qianxunke/book-ticket-common/proto/task"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"net/http"
)

func Init(client client.Client) *ApiService {
	return &ApiService{
		serviceClient: task.NewTaskService(basic.BookTicketService, client),
	}
}

type ApiService struct {
	serviceClient task.TaskService
}

func (api *ApiService) AddTask(c *gin.Context) {
	req := &task.In_AddTask{}
	req.TaskDetails = &task.TaskDetails{}
	if err := c.ShouldBindJSON(&req.TaskDetails); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &api_common.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	if len(c.Request.Header.Get("userId")) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, &api_common.Error{Code: http.StatusBadRequest, Message: "身份可疑"})
		return
	}

	req.TaskDetails.Task.UserId = c.Request.Header.Get("userId")
	rsp, _ := api.serviceClient.AddTask(context.TODO(), req)
	api_common.SrvResultDone(c, nil, &api_common.Error{Code: rsp.Error.Code, Message: rsp.Error.Message})
}

func (api *ApiService) GetTask(c *gin.Context) {
	req := &task.In_GetTaskInfo{}
	Id := c.Param("taskId")
	if len(Id) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, &api_common.Error{Code: http.StatusBadRequest, Message: "参数非法"})
		return
	}
	req.TaskId = Id
	rsp, _ := api.serviceClient.GetTaskInfo(context.TODO(), req)
	api_common.SrvResultDone(c, nil, &api_common.Error{Code: rsp.Error.Code, Message: rsp.Error.Message})
}
