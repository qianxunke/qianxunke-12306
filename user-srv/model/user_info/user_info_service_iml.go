package user_info

import (
	"book-user_srv/global"
	checkutil "book-user_srv/utils"
	"book-user_srv/utils/msm"
	"context"
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/plugins/db"
	auth "gitee.com/qianxunke/book-ticket-common/proto/auth"
	userInfoProto "gitee.com/qianxunke/book-ticket-common/proto/user"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func (s *userInfoServiceImp) DoneUserLogin(req *userInfoProto.InDoneUserLogin) (rsp *userInfoProto.OutDoneUserLogin) {
	rsp = &userInfoProto.OutDoneUserLogin{}
	if req.LoginType == 1 { //使用用户名。密码登陆
		loginByUserName(req, rsp)
	} else if req.LoginType == 2 { //验证码登陆
		loginByTelephone(req, rsp)
	} else {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusBadRequest,
			Message: "请求非法,未设置登陆方式！",
		}
	}
	if rsp.Error.Code != http.StatusOK {
		rsp.UserInf = nil
		return
	}
	rsp2, err := global.AuthClient.MakeAccessToken(context.TODO(), &auth.Request{
		UserId:   rsp.UserInf.UserId,
		UserName: rsp.UserInf.UserName,
	})
	if err != nil {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	_ = hideUserPrivate(rsp.UserInf)
	rsp.Token = rsp2.Token
	return
}

/**
  更新用户信息
*/
func (s *userInfoServiceImp) UpdateUserInfo(req *userInfoProto.InUpdateUserInfo) (rsp *userInfoProto.OutUpdateUserInfo) {
	rsp = &userInfoProto.OutUpdateUserInfo{}
	//这里只是修改普通信息
	updataData := map[string]interface{}{}
	if len(req.UserInf.Gender) > 0 {
		updataData["gender"] = req.UserInf.Gender
	}

	if len(req.UserInf.NikeName) > 0 {
		updataData["nike_name"] = req.UserInf.NikeName
	}
	if req.UserInf.IdentityCardType > 0 && len(req.UserInf.IdentityCardNo) > 0 {
		updataData["identity_card_type"] = req.UserInf.IdentityCardType
		updataData["identity_card_no"] = req.UserInf.IdentityCardNo
	}
	if len(updataData) <= 0 {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusBadRequest,
			Message: "没有任何数据需要修改！",
		}
		return
	}
	updataData["modified_time"] = time.Now().Format("2006-01-02 15:04:05")
	DB := db.MasterEngine()
	err := DB.Model(&userInfoProto.UserInf{}).Where("user_id =?", req.UserInf.UserId).Updates(updataData).Error
	if err != nil {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	rsp.Error = &userInfoProto.Error{
		Code:    http.StatusOK,
		Message: "修改成功！",
	}
	return
}

/**
获取用户列表
*/
func (s *userInfoServiceImp) GetUserInfoList(req *userInfoProto.InGetUserInfoList) (rsp *userInfoProto.OutGetUserInfoList) {
	//对参数鉴权
	if req.Limit == 0 {
		req.Limit = 10 //默认10个分页
	}
	if req.Limit > 1000 { //每一页数量
		req.Limit = 1000
	}
	if req.Pages <= 0 { //页数
		req.Pages = 1
	}
	DB := db.MasterEngine()
	rsp = &userInfoProto.OutGetUserInfoList{}
	var err error
	if len(req.SearchKey) == 0 {
		if len(req.StartTime) > 0 && len(req.EndTime) == 0 {
			err = DB.Model(&userInfoProto.UserInf{}).Where("register_time > ?", req.EndTime).Order("user_id desc").Count(&rsp.Total).Error
			if err == nil && rsp.Total > 0 {
				err = DB.Where("register_time > ? ", req.StartTime).Order("user_id desc").Offset((req.Pages - 1) * req.Limit).Limit(req.Limit).Find(&rsp.UserInfList).Error
			}
		} else if len(req.StartTime) == 0 && len(req.EndTime) > 0 {
			err = DB.Model(&userInfoProto.UserInf{}).Where("register_time < ? ", req.EndTime).Order("user_id desc").Count(&rsp.Total).Error
			if err == nil && rsp.Total > 0 {
				err = DB.Where("register_time < ? ", req.EndTime).Order("user_id desc").Offset((req.Pages - 1) * req.Limit).Limit(req.Limit).Find(&rsp.UserInfList).Error
			}
		} else if len(req.StartTime) > 0 && len(req.EndTime) > 0 {
			err = DB.Model(&userInfoProto.UserInf{}).Where("register_time  between ? and ?", req.StartTime, req.EndTime).Order("user_id desc").Count(&rsp.Total).Error
			if err == nil && rsp.Total > 0 {
				err = DB.Where("register_time  between ? and ?", req.StartTime, req.EndTime).Order("user_id desc").Offset((req.Pages - 1) * req.Limit).Limit(req.Limit).Find(&rsp.UserInfList).Error
			}
		} else {
			//先统计
			err = DB.Model(&userInfoProto.UserInf{}).Order("user_id desc").Count(&rsp.Total).Error
			if err == nil && rsp.Total > 0 {
				err = DB.Order("user_id desc").Offset((req.Pages - 1) * req.Limit).Limit(req.Limit).Find(&rsp.UserInfList).Error
			}
		}
	} else {
		key := "%" + req.SearchKey + "%"
		if len(req.StartTime) > 0 && len(req.EndTime) == 0 {
			err = DB.Model(&userInfoProto.UserInf{}).Where("(nike_name like ? or user_name like ? or mobile_phone like ?) and register_time > ? ", key, key, key, req.StartTime).Order("user_id desc").Count(&rsp.Total).Error
			if err == nil && rsp.Total > 0 {
				err = DB.Model(&userInfoProto.UserInf{}).Where("(nike_name like ? or user_name like ? or mobile_phone like ?) and register_time > ? ", key, key, key, req.StartTime).Order("user_id desc").Offset((req.Pages - 1) * req.Limit).Limit(req.Limit).Find(&rsp.UserInfList).Error
			}
		} else if len(req.StartTime) == 0 && len(req.EndTime) > 0 {
			err = DB.Model(&userInfoProto.UserInf{}).Where("(nike_name like ? or user_name like ? or mobile_phone like ?) and register_time < ? ", key, key, key, req.EndTime).Order("user_id desc").Count(&rsp.Total).Error
			if err == nil && rsp.Total > 0 {
				err = DB.Where("(nike_name like ? or user_name like ? or mobile_phone like ?) and register_time < ? ", key, key, key, req.EndTime).Order("user_id desc").Offset((req.Pages - 1) * req.Limit).Limit(req.Limit).Find(&rsp.UserInfList).Error
			}
		} else if len(req.StartTime) > 0 && len(req.EndTime) > 0 {
			err = DB.Model(&userInfoProto.UserInf{}).Where("(nike_name like ? or user_name like ? or mobile_phone like ?) and register_time between ? and ?", key, key, key, req.StartTime, req.EndTime).Order("user_id desc").Count(&rsp.Total).Error
			if err == nil && rsp.Total > 0 {
				err = DB.Where("(nike_name like ? or user_name like ? or mobile_phone like ?) and register_time between ? and ?", key, key, key, req.StartTime, req.EndTime).Order("user_id desc").Offset((req.Pages - 1) * req.Limit).Limit(req.Limit).Find(&rsp.UserInfList).Error
			}
		} else {
			err = DB.Model(&userInfoProto.UserInf{}).Where("1=1 and (nike_name like ? or user_name like ? or mobile_phone like ?)", key, key, key, req.StartTime).Order("user_id desc").Count(&rsp.Total).Error
			if err == nil && rsp.Total > 0 {
				err = DB.Where("1=1 and (nike_name like ? or user_name like ? or mobile_phone like ?)", key, key, key, req.StartTime).Order("user_id desc").Offset((req.Pages - 1) * req.Limit).Limit(req.Limit).Find(&rsp.UserInfList).Error
			}
		}
	}

	if err != nil {
		log.Printf("ERROR: %v", err)
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	var message string
	if rsp.Total > 0 && len(rsp.UserInfList) > 0 {
		for i := 0; i < len(rsp.UserInfList); i++ {
			_ = hideUserPrivate(rsp.UserInfList[i])
		}
		message = "查询成功！"
	} else {
		message = "没有数据了！"
	}
	//统计有多少条
	rsp.Error = &userInfoProto.Error{
		Code:    http.StatusOK,
		Message: message,
	}
	rsp.Limit = req.Limit
	rsp.Pages = req.Pages
	return
}

/**
  获取用户信息
*/
func (s *userInfoServiceImp) GetUserInfo(req *userInfoProto.InGetUserInfo) (rsp *userInfoProto.OutGetUserInfo) {
	rsp = &userInfoProto.OutGetUserInfo{}
	rsp.UserInf = &userInfoProto.UserInf{}
	//获取某个用户的信息
	DB := db.MasterEngine()
	err := DB.Where("user_id = ?", req.UserId).First(&rsp.UserInf).Error
	if err != nil && rsp.UserInf.UserId == 0 {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusBadRequest,
			Message: "没有找到任何数据！",
		}
		return
	}
	rsp.Roles = []string{"admin"}
	_ = hideUserPrivate(rsp.UserInf)
	rsp.Error = &userInfoProto.Error{
		Code:    http.StatusOK,
		Message: "OK",
	}
	return
}

/**
  注册
*/
func (s *userInfoServiceImp) DoeUserRegister(req *userInfoProto.InDoneUserRegister) (rsp *userInfoProto.OutDoneUserRegister) {
	rsp = &userInfoProto.OutDoneUserRegister{}
	userInf := &userInfoProto.UserInf{}
	DB := db.MasterEngine()
	err := DB.Where(" mobile_phone = ? or user_name= ?", req.Userinf.MobilePhone, req.Userinf.UserName).First(&userInf).Error
	if err == nil || userInf.UserId > 0 {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusBadRequest,
			Message: "用户已存在",
		}
		return
	}
	if len(req.VerificationCode) == 0 {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusBadRequest,
			Message: "验证码为空",
		}
		return
	}
	err = verificationTelphone(req.Userinf.MobilePhone)
	if err != nil {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return
	}
	md5str, err := getMd5Password(req.Userinf.Password)
	if err != nil {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	req.Userinf.Password = md5str
	req.Userinf.UserStats = 1
	req.Userinf.RegisterTime = time.Now().Format("2006-01-02 15:04:05")
	req.Userinf.UserPoint = 0
	req.Userinf.ModifiedTime = req.Userinf.RegisterTime
	err = DB.Create(req.Userinf).Error
	if err != nil {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusInternalServerError,
			Message: "创建用户失败",
		}
		return
	}
	rsp.UserInf = &userInfoProto.UserInf{}
	err = DB.Where(" mobile_phone = ?", req.Userinf.MobilePhone).First(&rsp.UserInf).Error
	if err != nil || rsp.UserInf.UserId <= 0 {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusInternalServerError,
			Message: "查询用户失败",
		}
		return
	}
	rsp2, err := global.AuthClient.MakeAccessToken(context.TODO(), &auth.Request{
		UserId:   rsp.UserInf.UserId,
		UserName: rsp.UserInf.UserName,
	})
	if err != nil {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusInternalServerError,
			Message: "token 生成出错",
		}
		return
	}
	rsp.Token = rsp2.Token
	rsp.Error = &userInfoProto.Error{
		Code:    http.StatusOK,
		Message: "注册成功",
	}
	_ = hideUserPrivate(rsp.UserInf)
	return
}

/**
 *验证码发送
 */
func (s *userInfoServiceImp) GetVerificationCode(req *userInfoProto.InGetVerificationCode) (rsp *userInfoProto.OutGetVerificationCode) {
	rsp = &userInfoProto.OutGetVerificationCode{}
	if req.Telephone == "" {
		rsp.Error = &userInfoProto.Error{
			Message: "电话号码为空",
			Code:    http.StatusBadRequest,
		}
		return
	}

	if checkutil.ValiTephone(req.Telephone) {
		rsp.Error = &userInfoProto.Error{
			Message: "电话号码不合法",
			Code:    http.StatusBadRequest,
		}
		return
	}

	//判断该手机验证码是否还有效
	err := verificationTelphone(req.Telephone)
	if err == nil {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusBadRequest,
			Message: "验证码已发，请耐心等等，或请一分钟后再次请求！",
		}
		return
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	int6, err := strconv.ParseInt(vcode, 10, 64)
	if err != nil {
		rsp.Error = &userInfoProto.Error{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
		return
	}
	err = msm.SendRegisterMsm(int6, req.Telephone, global.RedisClient)
	if err != nil {
		rsp.Error = &userInfoProto.Error{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
		return
	}
	rsp.Error = &userInfoProto.Error{
		Code:    http.StatusOK,
		Message: "ok",
	}
	return
}

//获取客户联系人
func (s *userInfoServiceImp) GetUserPassengerList(req *userInfoProto.In_GetUserPassengerList) (rsp *userInfoProto.Out_GetUserPassengerList) {
	//从数据库获取，没有的话想办法获取
	DB := db.MasterEngine()
	rsp = &userInfoProto.Out_GetUserPassengerList{}
	err := DB.Model(&userInfoProto.Passenger{}).Where("user_id", req.UserId).Find(rsp.PassengerList).Error
	if err != nil {
		log.Printf("【GetUserPassengerList】error : %s", err.Error())
		rsp.Error = &userInfoProto.Error{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
		return
	}
	return
}

//修改用户联系人
func (s *userInfoServiceImp) UpdateUserPassenger(req *userInfoProto.In_UpdateUserPassenger) (rsp *userInfoProto.Out_UpdateUserPassenger) {
	//这里异步即可,不需要
	return
}
