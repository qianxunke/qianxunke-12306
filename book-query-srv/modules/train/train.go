package train

import (
	"book-query-srv/modules/train/dao"
	"book-query-srv/modules/train/service"
)

func Init() {
	dao.Init()
	service.Init()
}
