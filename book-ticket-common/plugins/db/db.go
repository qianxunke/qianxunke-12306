package db

import (
	"gitee.com/qianxunke/surprise-shop-common/basic/config"
	"github.com/jinzhu/gorm"
	"github.com/kataras/golog"
	"github.com/micro/go-micro/util/log"
	"sync"
)

var (
	masterEngine *gorm.DB //主数据库
	slaveEngine  *gorm.DB //从数据库
	lock         sync.Mutex
)

func init() {
	//basic.Register
}

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
		golog.Errorf("@@@ 数据库 master 节点连接异常挂掉!! %s", err)
		createEngine(true)
	}
	return masterEngine
}

// 从库，单例
func SlaveEngine() *gorm.DB {
	if slaveEngine != nil {
		goto EXIST
	}
	lock.Lock()
	defer lock.Unlock()

	if slaveEngine != nil {
		goto EXIST
	}

	createEngine(false)
	return slaveEngine

EXIST:
	var err = slaveEngine.DB().Ping()
	if err != nil {
		golog.Errorf("@@@ 数据库 slave 节点连接异常挂掉!! %s", err)
		createEngine(false)
	}
	return slaveEngine
}

func createEngine(isMaster bool) {
	c := config.C()
	cfg := &db{}
	err := c.App("db", cfg)
	if err != nil {
		log.Logf("[initMysql] %s", err)
	}

	if !cfg.Mysql.Enable {
		log.Logf("[initMysql] 未启用Mysql")
		return
	}

	engine, err := gorm.Open("mysql", cfg.Mysql.URL)
	if err != nil {
		golog.Fatalf("@@@ 初始化数据库连接失败!! %s", err)
		return
	}
	//是否启用日志记录器，将会在控制台打印sql
	engine.LogMode(true)
	if cfg.Mysql.MaxIdleConnection > 0 {
		engine.DB().SetMaxIdleConns(cfg.Mysql.MaxIdleConnection)
	}
	if cfg.Mysql.MaxOpenConnection > 0 {
		engine.DB().SetMaxOpenConns(cfg.Mysql.MaxOpenConnection)
	}
	if isMaster {
		masterEngine = engine
	} else {
		slaveEngine = engine
	}
}
