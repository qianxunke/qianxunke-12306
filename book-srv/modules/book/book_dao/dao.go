package book_dao

import (
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/plugins/db"
	"gitee.com/qianxunke/book-ticket-common/proto/task"
	"sync"
)

//任务数据库操作

var (
	dao *taskDaoIml
	m   sync.Mutex
)

type taskDaoIml struct {
}

type TaskDao interface {
	FindById(id string) (task *task.TaskDetails, err error)

	Insert(task *task.TaskDetails) (err error)

	SimpleQuery(limit int64, pages int64, status int64, key string, startTime string, endTime string, order string) (rsp *task.Out_GetTaskInfoList, err error)

	Delete(ids []int64) (err error)

	Update(task *task.TaskDetails) (err error)

	TicketQuery(limit int64, pages int64, status int64) (rsp []task.TaskDetails, err error)
}

func GetDao() (TaskDao, error) {
	if dao == nil {
		return nil, fmt.Errorf("[GetService] GetService 未初始化")
	}
	return dao, nil
}

func Init() {
	m.Lock()
	defer m.Unlock()
	if dao != nil {
		return
	}
	// 检查模型是否存在
	master := db.MasterEngine()
	if !master.HasTable(&task.Task{}) {
		master.CreateTable(&task.Task{})
	}
	if !master.HasTable(&task.TaskPassenger{}) {
		master.CreateTable(&task.TaskPassenger{})
	}
	dao = &taskDaoIml{}
}
