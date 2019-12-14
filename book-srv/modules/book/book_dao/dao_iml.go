package book_dao

import (
	"errors"
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/plugins/db"
	"gitee.com/qianxunke/book-ticket-common/proto/task"
	"log"
)

func (dao *taskDaoIml) FindById(taskId string) (product *task.TaskDetails, err error) {
	product = &task.TaskDetails{}
	DB := db.MasterEngine()
	err = DB.Where("task_id = ?", taskId).First(&product.Task).Error
	err = DB.Where("task_id = ?", taskId).Find(&product.TaskPassenger).Error
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

	err = DB.Create(&product).Error
	if err != nil {
		DB.Rollback()
		return
	}
	for _, item := range product.TaskPassenger {
		err = DB.Save(item).Error
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
	err = tx.Model(&task.Task{}).Where("task_id = ?", ta.Task.TaskId).Updates(ta).Error
	if err != nil {
		log.Printf("[Update] error %s", err.Error())
		tx.Rollback()
		return
	}
	err = tx.Raw("delete from task_passengers where task_id = ?", ta.Task.TaskId).Error
	if err != nil {
		log.Printf("[Update] error %s", err.Error())
		tx.Rollback()
		return
	}
	for _, item := range ta.TaskPassenger {
		err = tx.Save(item).Error
		if err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return
}
