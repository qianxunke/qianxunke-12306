package book_dao

import (
	"errors"
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/basic/common/uuid"
	"gitee.com/qianxunke/book-ticket-common/plugins/db"
	"gitee.com/qianxunke/book-ticket-common/proto/task"
	"log"
)

func (dao *taskDaoIml) FindById(taskId string) (product *task.TaskDetails, err error) {
	product = &task.TaskDetails{}
	product.Task = &task.Task{}
	product.TaskPassenger = []*task.TaskPassenger{}
	DB := db.MasterEngine()
	err = DB.Model(&task.Task{}).Where("task_id = ?", taskId).First(&product.Task).Error
	if err != nil {
		log.Println("FindById 1:" + err.Error())
	}
	err = DB.Where("task_id = ?", taskId).Find(&product.TaskPassenger).Error
	if err != nil {
		log.Println("FindById 2:" + err.Error())
	}
	return
}

func (dao *taskDaoIml) Insert(product *task.TaskDetails) (err error) {
	DB := db.MasterEngine()
	DB.Begin()
	defer func() {
		if re := recover(); re != nil {
			DB.Rollback()
			err = errors.New(fmt.Sprintf("%v", re))
			log.Printf("[Update] error %s", err.Error())
		}
	}()
	product.Task.TaskId = uuid.GetUuid()
	err = DB.Create(&product.Task).Error
	if err != nil {
		DB.Rollback()
		return
	}
	for _, item := range product.TaskPassenger {
		item.Id = uuid.GetUuid()
		item.TaskId = product.Task.TaskId
		err = DB.Create(item).Error
		if err != nil {
			DB.Rollback()
			return
		}
	}
	return
}

func (dao *taskDaoIml) SimpleQuery(limit int64, pages int64, status int64, key string, startTime string, endTime string, order string) (rsp *task.Out_GetTaskInfoList, err error) {
	DB := db.MasterEngine()
	rsp = &task.Out_GetTaskInfoList{}
	offset := (pages - 1) * limit
	if len(key) == 0 {
		if len(startTime) > 0 && len(endTime) == 0 {
			err = DB.Model(&task.Task{}).Where("status = ? and created_time > ?", status, endTime).Order(order).Count(&rsp.Total).Error
			if err == nil && rsp.Total > 0 {
				err = DB.Where("status = ? and created_time > ? ", status, startTime).Order(order).Offset(offset).Limit(limit).Find(&rsp.TaskDetailsList).Error
			}
		} else if len(startTime) == 0 && len(endTime) > 0 {
			err = DB.Model(&task.Task{}).Where("status = ? and created_time < ? ", status, endTime).Order(order).Count(&rsp.Total).Error
			if err == nil && rsp.Total > 0 {
				err = DB.Where("status = ? and created_time < ? ", status, endTime).Order(order).Offset(offset).Limit(limit).Find(&rsp.TaskDetailsList).Error
			}
		} else if len(startTime) > 0 && len(endTime) > 0 {
			err = DB.Model(&task.Task{}).Where("status = ? and created_time  between ? and ?", status, startTime, endTime).Order(order).Count(&rsp.Total).Error
			if err == nil && rsp.Total > 0 {
				err = DB.Where("status = ? and created_time  between ? and ?", status, startTime, endTime).Order(order).Offset(offset).Limit(limit).Find(&rsp.TaskDetailsList).Error
			}
		} else {
			//先统计
			err = DB.Model(&task.Task{}).Where("status = ? ", status).Order(order).Count(&rsp.Total).Error
			if err == nil && rsp.Total > 0 {
				err = DB.Model(&task.Task{}).Where("status = ? ", status).Order(order).Offset(offset).Limit(limit).Find(&rsp.TaskDetailsList).Error
			}
		}
	} else {
		searchKey := "%" + key + "%"
		if len(startTime) > 0 && len(endTime) == 0 {
			err = DB.Model(&task.Task{}).Where("status = ? and (name like ? ) and created_time > ? ", status, searchKey, startTime).Order(order).Count(&rsp.Total).Error
			if err == nil && rsp.Total > 0 {
				err = DB.Model(&task.Task{}).Where("status = ? and (name like ?) and created_time > ? ", status, searchKey, startTime).Order(order).Offset(offset).Limit(limit).Find(&rsp.TaskDetailsList).Error
			}
		} else if len(startTime) == 0 && len(endTime) > 0 {
			err = DB.Model(&task.Task{}).Where("status = ? and (name like ?) and created_time < ? ", status, searchKey, endTime).Order(order).Count(&rsp.Total).Error
			if err == nil && rsp.Total > 0 {
				err = DB.Where("status = ? and (name like ?) and created_time < ? ", status, searchKey, endTime).Order(order).Offset(offset).Limit(limit).Find(&rsp.TaskDetailsList).Error
			}
		} else if len(startTime) > 0 && len(endTime) > 0 {
			err = DB.Model(&task.Task{}).Where("status = ? and (name like ?) and created_time between ? and ?", status, searchKey, startTime, endTime).Order(order).Count(&rsp.Total).Error
			if err == nil && rsp.Total > 0 {
				err = DB.Where("status = ? and (name like ?) and created_time between ? and ?", status, searchKey, startTime, endTime).Order(order).Offset(offset).Limit(limit).Find(&rsp.TaskDetailsList).Error
			}
		} else {
			err = DB.Model(&task.Task{}).Where("status = ? and name like ?", status, searchKey).Order(order).Count(&rsp.Total).Error
			if err == nil && rsp.Total > 0 {
				err = DB.Where("status = ? and name like ?", status, searchKey).Order(order).Offset(offset).Limit(limit).Find(&rsp.TaskDetailsList).Error
			}
		}
	}
	return
}

func (dao *taskDaoIml) Delete(ids []int64) (err error) {
	if len(ids) == 0 {
		return
	}
	DB := db.MasterEngine()
	DB.Begin()
	defer func() {
		if err != nil {
			DB.Rollback()
		}
	}()
	for i := 0; i < len(ids); i++ {
		err = DB.Where("task_id = ?", ids[i]).Delete(&task.TaskPassenger{}).Error
		if err != nil {
			break
		}
		err = DB.Where("task_id = ?", ids[i]).Delete(&task.Task{}).Error
		if err != nil {
			break
		}
	}
	if err == nil {
		DB.Commit()
	}
	return
}

func (dao *taskDaoIml) Update(ta *task.TaskDetails) (err error) {
	DB := db.MasterEngine()
	tx := DB.Begin()
	defer func() {
		if re := recover(); re != nil {
			err = errors.New(fmt.Sprintf("%v", re))
			log.Printf("[Update] error %s", err.Error())
			tx.Rollback()
			return
		}

	}()
	err = tx.Model(&task.Task{}).Where("task_id = ?", ta.Task.TaskId).Update(ta.Task).Error
	if err != nil {
		log.Printf("[Update] error %s", err.Error())
		tx.Rollback()
		return
	}
	err = tx.Delete(task.TaskPassenger{}, "task_id =?", ta.Task.TaskId).Error
	if err != nil {
		log.Printf("[Update] error %s", err.Error())
		tx.Rollback()
		return
	}
	for _, item := range ta.TaskPassenger {
		item.Id = uuid.GetUuid()
		item.TaskId = ta.Task.TaskId
		err = tx.Create(item).Error
		if err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return
}

func (dao *taskDaoIml) TicketQuery(limit int64, pages int64, status int64) (rsp []task.TaskDetails, err error) {
	DB := db.MasterEngine()
	total := 0
	offset := (pages - 1) * limit
	//先查询task
	err = DB.Model(&task.Task{}).Where("status = ?", status).Count(&total).Error
	if err != nil {
		return nil, err
	}
	if total <= 0 {
		return make([]task.TaskDetails, 0), err
	}
	var tasks []*task.Task
	err = DB.Where("status =  ?", status).Offset(offset).Limit(limit).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	for _, item := range tasks {
		var presenters []*task.TaskPassenger
		err = DB.Where("task_id =  ?", item.TaskId).Find(&presenters).Error
		if err != nil {
			return nil, err
		}
		rsp = append(rsp, task.TaskDetails{Task: item, TaskPassenger: presenters})
	}
	return
}
