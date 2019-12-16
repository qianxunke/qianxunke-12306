package user_info

import (
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/plugins/db"
	userInfoProto "gitee.com/qianxunke/book-ticket-common/proto/user"
	"sync"
)

var (
	s *userInfoServiceImp
	m sync.Mutex
)

//service 服务
type userInfoServiceImp struct {
}

//用户服务类
type UserService interface {
	//用户登陆
	DoneUserLogin(req *userInfoProto.InDoneUserLogin) (rsp *userInfoProto.OutDoneUserLogin)
	//修改用户信息
	UpdateUserInfo(req *userInfoProto.InUpdateUserInfo) (rsp *userInfoProto.OutUpdateUserInfo)
	//获取用户列表
	GetUserInfoList(req *userInfoProto.InGetUserInfoList) (rsp *userInfoProto.OutGetUserInfoList)
	//获取用户详细信息
	GetUserInfo(req *userInfoProto.InGetUserInfo) (rsp *userInfoProto.OutGetUserInfo)
	//用户注册
	DoeUserRegister(req *userInfoProto.InDoneUserRegister) (rsp *userInfoProto.OutDoneUserRegister)
	//获取验证码
	GetVerificationCode(req *userInfoProto.InGetVerificationCode) (rsp *userInfoProto.OutGetVerificationCode)
	//获取客户联系人
	GetUserPassengerList(req *userInfoProto.In_GetUserPassengerList) (rsp *userInfoProto.Out_GetUserPassengerList)
	//修改用户联系人
	UpdateUserPassenger(req *userInfoProto.In_UpdateUserPassenger) (rsp *userInfoProto.Out_UpdateUserPassenger)
	//登陆12306客户
	Login12306 (req *userInfoProto.In_Login12306)  (rsp *userInfoProto.Out_Login12306)
}

func GetService() (UserService, error) {
	if s == nil {
		return nil, fmt.Errorf("[GetService] GetService 未初始化")
	}
	return s, nil
}

//初始化用户服务层
func Init() {
	m.Lock()
	defer m.Unlock()

	if s != nil {
		return
	}
	master := db.MasterEngine()
	if !master.HasTable(&userInfoProto.UserInf{}){
		master.CreateTable(&userInfoProto.UserInf{})
	}
	if !master.HasTable(&userInfoProto.Passenger{}){
		master.CreateTable(&userInfoProto.Passenger{})
	}
	s = &userInfoServiceImp{}
}
