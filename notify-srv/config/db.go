package config

import (
	"github.com/jinzhu/gorm"
	"log"
	"sync"
)

/**
  解析数据库配置工具
*/

var (
	masterEngine *gorm.DB //主数据库
	lock         sync.Mutex
)

//配置数据库主库
func MasterEngine() *gorm.DB {

	if masterEngine != nil {
		goto EXiST
	}
	//锁住
	lock.Lock()
	defer lock.Unlock()
	if masterEngine != nil {
		goto EXiST
	}
	createEngine(true)
	return masterEngine

EXiST:
	var err = masterEngine.DB().Ping()
	if err != nil {
		log.Printf("@@@ 数据库 master 节点连接异常挂掉!! %s", err)
		createEngine(true)
	}
	return masterEngine
}

func createEngine(isMaster bool) {

	engine, err := gorm.Open("sqlite3", "qianxunke_ticket_citys.db")
	if err != nil {
		log.Printf("@@@ 初始化数据库连接失败!! %s", err)
		return
	}
	//是否启用日志记录器，将会在控制台打印sql
	engine.LogMode(true)

	engine.DB().SetMaxIdleConns(50)

	engine.DB().SetMaxOpenConns(50)

	masterEngine = engine

}
