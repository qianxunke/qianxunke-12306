package user_info

import (
	userIml "book-user_srv/model/user_info"
	"context"
	userProto "gitee.com/qianxunke/book-ticket-common/proto/user"
	"github.com/micro/go-micro/util/log"
)

type Handler struct{}

var (
	userService userIml.UserService
)

func Init() {
	_userService, err := userIml.GetService()
	if err != nil {
		log.Fatal("[Init] 初始化Handler错误")
		return
	}
	userService = _userService
}

func (e *Handler) DoneUserLogin(ctx context.Context, req *userProto.InDoneUserLogin, rsp *userProto.OutDoneUserLogin) error {
	log.Log("Received CostomerInfo.DoneUserLogin request")
	response := userService.DoneUserLogin(req)
	rsp.Error = response.Error
	rsp.UserInf = response.UserInf
	rsp.Token = response.Token
	return nil
}

func (e *Handler) GetUserInfo(ctx context.Context, req *userProto.InGetUserInfo, rsp *userProto.OutGetUserInfo) error {
	log.Log("Received UserInfo.GetUserInfo request")
	response := userService.GetUserInfo(req)
	rsp.Error = response.Error
	rsp.UserInf = response.UserInf
	rsp.Roles = response.Roles
	return nil
}

func (e *Handler) GetUserInfoList(ctx context.Context, req *userProto.InGetUserInfoList, rsp *userProto.OutGetUserInfoList) error {
	log.Log("Received UserInfo.GetUserInfoList request")
	response := userService.GetUserInfoList(req)
	rsp.Error = response.Error
	rsp.UserInfList = response.UserInfList
	rsp.Pages = response.Pages
	rsp.Limit = response.Limit
	rsp.Total = response.Total
	return nil
}

func (e *Handler) UpdateUserInfo(ctx context.Context, req *userProto.InUpdateUserInfo, rsp *userProto.OutUpdateUserInfo) error {
	log.Logf("Received UserInfo.UpdateUserInfo request : %v", req)
	response := userService.UpdateUserInfo(req)
	rsp.Error = response.Error
	rsp.UserInf = response.UserInf
	return nil
}

func (e *Handler) DoneUserRegister(ctx context.Context, req *userProto.InDoneUserRegister, rsp *userProto.OutDoneUserRegister) error {
	log.Log("Received UserInfo.DoeUserRegister request")
	response := userService.DoeUserRegister(req)
	rsp.Error = response.Error
	rsp.UserInf = response.UserInf
	rsp.Token = response.Token
	return nil
}

func (e *Handler) GetVerificationCode(ctx context.Context, req *userProto.InGetVerificationCode, rsp *userProto.OutGetVerificationCode) error {
	log.Log("Received UserInfo.GetVerificationCode request")
	response := userService.GetVerificationCode(req)
	rsp.Error = response.Error
	return nil
}

func (e *Handler) GetUserPassengerList(ctx context.Context, req *userProto.In_GetUserPassengerList, rsp *userProto.Out_GetUserPassengerList) error {
	log.Log("Received UserInfo.GetUserPassengerList request")
	response := userService.GetUserPassengerList(req)
	rsp.Error = response.Error
	rsp.PassengerList = response.PassengerList
	return nil
}

func (e *Handler) UpdateUserPassenger(ctx context.Context, req *userProto.In_UpdateUserPassenger, rsp *userProto.Out_UpdateUserPassenger) error {
	log.Log("Received UserInfo.UpdateUserPassenger request")
	response := userService.UpdateUserPassenger(req)
	rsp.Error = response.Error
	return nil
}
func (e *Handler) Login12306(ctx context.Context, req *userProto.In_Login12306, rsp *userProto.Out_Login12306) error {
	log.Log("Received UserInfo.Login12306 request")
	response := userService.Login12306(req)
	rsp.Error = response.Error
	return nil
}
