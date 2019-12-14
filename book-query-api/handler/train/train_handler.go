package train

import (
	"context"
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/basic"
	"gitee.com/qianxunke/book-ticket-common/basic/api_common"
	"gitee.com/qianxunke/book-ticket-common/proto/ticket"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"time"
)

func Init(client client.Client) *ApiService {
	return &ApiService{
		serviceClient: ticket.NewTrainService(basic.InventoryService, client),
	}
}

type ApiService struct {
	serviceClient ticket.TrainService
}

func (as *ApiService) GetTrainInfoList(c *gin.Context) {
	reqParameter := &ticket.In_GetTrainInfoList{}
	reqParameter.TrainDate = c.DefaultQuery("train_date", fmt.Sprintf("%d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day()))
	reqParameter.FindFrom = c.DefaultQuery("find_from", "北京")
	reqParameter.FindTo = c.DefaultQuery("find_to", "上海")
	reqParameter.PurposeCodes = c.DefaultQuery("purpose_codes", "ADULT")
	rsp, _ := as.serviceClient.GetTrainInfoList(context.TODO(), reqParameter)
	//返回结果
	api_common.SrvResultListDone(c, rsp.TrainList, rsp.Limit, rsp.Pages, rsp.Total, &api_common.Error{Code: rsp.Error.Code, Message: rsp.Error.Message})

}
