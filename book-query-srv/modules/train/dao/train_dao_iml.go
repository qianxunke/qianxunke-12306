package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/plugins/db"
	ticketProto "gitee.com/qianxunke/book-ticket-common/proto/ticket"
	"log"
)

func (dao *ticketDaoIml) FindById(secretStr string) (train *ticketProto.Train, err error) {
	train = &ticketProto.Train{}
	DB := db.MasterEngine()
	err = DB.Where("secret_str = ?", secretStr).First(&train).Error
	return
}

func (dao *ticketDaoIml) Insert(train []*ticketProto.Train) (err error) {
	defer func() {
		if re := recover(); re != nil {
			err = errors.New(fmt.Sprintf("%v", re))
			log.Println("[Insert] error : " + err.Error())
		}
	}()
	jsons, errs := json.Marshal(train) //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
	}
	err = redisClient.Set(train[0].TrainDate+"_"+train[0].FindFrom+"_"+train[0].FindTo+"_"+train[0].PurposeCodes, string(jsons), tokenExpiredDate).Err()
	if err != nil {
		log.Println("[Insert] error : " + err.Error())
	}
	return
}

func (dao *ticketDaoIml) SimpleQuery(findFrom string, findTo string, trainDate string, purposeCodes string) (value string, err error) {
	defer func() {
		if re := recover(); re != nil {
			err = errors.New(fmt.Sprintf("%v", re))
			log.Println("[SimpleQuery] error : " + err.Error())
		}
	}()
	value, err = redisClient.Get(trainDate + "_" + findFrom + "_" + findTo + "_" + purposeCodes).Result()
	if err != nil {
		log.Println("[SimpleQuery] error : " + err.Error())
	}
	return
}

func (dao *ticketDaoIml) Delete(ids []string) (err error) {
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
		err = DB.Where("secret_str = ?", ids[i]).Delete(&ticketProto.Train{}).Error
		if err != nil {
			break
		}
	}
	if err != nil {
		DB.Commit()
	}
	return
}

func (dao *ticketDaoIml) Update(product *ticketProto.Train) (err error) {
	DB := db.MasterEngine()
	err = DB.Model(&ticketProto.Train{}).Where("secret_str = ?", product.SecretStr).Updates(&product).Error
	return
}
