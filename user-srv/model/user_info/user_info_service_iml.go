package user_info

import (
	"book-user_srv/global"
	"book-user_srv/utils"
	"context"
	"fmt"
	"gitee.com/qianxunke/book-ticket-common/basic/common/uuid"
	"gitee.com/qianxunke/book-ticket-common/notice/sms"
	"gitee.com/qianxunke/book-ticket-common/plugins/db"
	"gitee.com/qianxunke/book-ticket-common/proto/auth"
	userInfoProto "gitee.com/qianxunke/book-ticket-common/proto/user"
	bookBean "gitee.com/qianxunke/book-ticket-common/ticket/book/bean"
	"gitee.com/qianxunke/book-ticket-common/ticket/book/boo_core"
	"gitee.com/qianxunke/book-ticket-common/ticket/login"
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
	//	_ = hideUserPrivate(rsp.UserInf)
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
	if err != nil || len(rsp.UserInf.UserId) == 0 {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusBadRequest,
			Message: "没有找到任何数据！",
		}
		return
	}
	rsp.Roles = []string{"admin"}
	//_ = hideUserPrivate(rsp.UserInf)
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
	err := DB.Where(" mobile_phone = ?", req.Userinf.MobilePhone).First(&userInf).Error
	if err == nil || len(userInf.UserId) > 0 {
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
	err = verificationTelphone(req.VerificationCode, req.Userinf.MobilePhone)
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
	req.Userinf.UserId = uuid.GetUuid()
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
	if err != nil || len(rsp.UserInf.UserId) <= 0 {
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
	//	_ = hideUserPrivate(rsp.UserInf)
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
	if req.Code == "login" {
		DB := db.MasterEngine()
		//判断该手机号是否注册
		user := &userInfoProto.UserInf{}
		err := DB.Where(" mobile_phone = ?", req.Telephone).First(&user).Error
		if err != nil {
			rsp.Error = &userInfoProto.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			return
		}

		if len(user.UserId) == 0 {
			rsp.Error = &userInfoProto.Error{
				Code:    http.StatusBadRequest,
				Message: "该手机号未注册，请先注册",
			}
			return
		}
	}

	//判断该手机验证码是否还有效
	err := verificationCodeIsOk(req.Telephone)
	if err == nil {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusBadRequest,
			Message: "验证码已发，请耐心等等，或请3分钟后再次请求！",
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
	err = sms.SendRegisterMsm(int6, req.Telephone, global.RedisClient)
	if err != nil {
		rsp.Error = &userInfoProto.Error{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
		return
	}
	rsp.Error = &userInfoProto.Error{
		Code:    http.StatusOK,
		Message: "验证码已发送...",
	}
	return
}

//获取客户联系人
func (s *userInfoServiceImp) GetUserPassengerList(req *userInfoProto.In_GetUserPassengerList) (rsp *userInfoProto.Out_GetUserPassengerList) {
	//从数据库获取，没有的话想办法获取
	DB := db.MasterEngine()
	rsp = &userInfoProto.Out_GetUserPassengerList{}
	defer func() {
		if re := recover(); re != nil {
			rsp.Error = &userInfoProto.Error{
				Message: fmt.Sprintf("%v", re),
				Code:    http.StatusInternalServerError,
			}
		}
	}()
	err := DB.Model(&userInfoProto.Passenger{}).Where("user_id = ?", req.UserId).Find(&rsp.PassengerList).Error
	if err != nil {
		log.Printf("【GetUserPassengerList】error : %s", err.Error())
		rsp.Error = &userInfoProto.Error{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
		return
	} else {
		rsp.Error = &userInfoProto.Error{
			Message: "ok",
			Code:    http.StatusOK,
		}
	}
	return
}

//修改用户联系人
func (s *userInfoServiceImp) UpdateUserPassenger(req *userInfoProto.In_UpdateUserPassenger) (rsp *userInfoProto.Out_UpdateUserPassenger) {
	//这里异步即可,不需要
	return
}

func (s *userInfoServiceImp) Login12306(req *userInfoProto.In_Login12306) (rsp *userInfoProto.Out_Login12306) {
	//判断用户在不在
	rsp = &userInfoProto.Out_Login12306{}
	DB := db.MasterEngine()
	userInf := userInfoProto.UserInf{}
	err := DB.Where(" user_id = ?", req.UserId).First(&userInf).Error
	if err != nil {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	if len(userInf.UserId) == 0 {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusBadRequest,
			Message: "用户不存在",
		}
		return
	}
	userInf.TranUserAccount = req.TranUserAccount
	userInf.TranUserPwd = req.TranUserPwd
	//开始登陆12306
	loginErrNum := 0
	var loginResult *login.LoginResult
	for true {
		if loginErrNum > 3 {
			break
		}
		t, err := login.LoginAndCheckToken(userInf)
		if err != nil {
			loginErrNum++
		} else {
			if t != nil {
				loginResult = t
				break
			}
		}
		time.Sleep(time.Second * 3)
	}
	if loginErrNum <= 3 && loginResult != nil {
		log.Println("登陆成功")
	} else {
		log.Println("登陆失败")
		return
	}

	userInf.TranUserName = loginResult.Username
	//修改用户12306账号和密码
	DB.Model(&userInf).Where("user_id = ?", userInf.UserId).Updates(userInfoProto.UserInf{TranUserName: userInf.TranUserName, TranUserAccount: userInf.TranUserAccount, TranUserPwd: userInf.TranUserPwd})

	//获取用户的联系人
	presenterErrNum := 0
	var presenterList []bookBean.Normal_passengers
	for true {
		if presenterErrNum > 3 {
			break
		}
		subToken, _, err := boo_core.GetInitDc(loginResult.Conversat)
		if err != nil {
			presenterErrNum++
			continue
		}
		ok, pData, err := boo_core.GetPassenger(http.MethodPost, loginResult.Conversat, subToken)
		if ok {
			if len(pData) > 0 {
				presenterList = pData
				break
			}
		}
		time.Sleep(time.Second * 1)
	}

	if presenterErrNum > 3 {
		rsp.Error = &userInfoProto.Error{
			Code:    http.StatusInternalServerError,
			Message: "登陆成功,但获取联系人失败",
		}
		return
	}
	//转化联系人
	temp := make([]userInfoProto.Passenger, len(presenterList))
	for i, item := range presenterList {
		temp[i].Id = uuid.GetUuid()
		temp[i].UserId = req.UserId
		temp[i].AllEncStr = item.AllEncStr
		temp[i].Address = item.Address
		temp[i].BornDate = item.Born_date
		temp[i].CountryCode = item.Country_code
		temp[i].Email = item.Email
		temp[i].FirstLetter = item.First_letter
		temp[i].GatValidDateEnd = item.Gat_valid_date_end
		temp[i].GatValidDateStart = item.Gat_valid_date_start
		temp[i].GatVersion = item.Gat_version
		temp[i].IndexId = item.Index_id
		temp[i].IsAdult = item.IsAdult
		temp[i].IsOldThan60 = item.IsOldThan60
		temp[i].IsYongThan10 = item.IsYongThan10
		temp[i].IsYongThan14 = item.IsYongThan14
		temp[i].MobileNo = item.Mobile_no
		temp[i].PassengerFlag = item.Passenger_flag
		temp[i].PassengerIdNo = item.Passenger_id_no
		temp[i].PassengerIdTypeCode = item.Passenger_id_type_code
		temp[i].PassengerIdTypeName = item.Passenger_id_type_name
		temp[i].PassengerName = item.Passenger_name
		temp[i].PassengerTypeName = item.Passenger_type_name
		temp[i].PhoneNo = item.Phone_no
		temp[i].Postalcode = item.Postalcode
		temp[i].SexCode = item.Sex_code
		temp[i].SexName = item.Sex_name
		temp[i].TotalTimes = item.Total_times
	}

	//删除之前的信息
	DB.Model(&userInfoProto.Passenger{}).Where("user_id = ?", userInf.UserId).Delete(&userInfoProto.Passenger{})

	//写进去
	for _, item := range temp {
		DB.Create(&item)
	}

	rsp.Error = &userInfoProto.Error{
		Code:    http.StatusOK,
		Message: "登陆成功",
	}
	return
}

//获取版本信息
func (s *userInfoServiceImp) GetUpdateInfo(req *userInfoProto.In_UpdateInfo) (rsp *userInfoProto.Out_UpdateInfo) {
	rsp = &userInfoProto.Out_UpdateInfo{}
	rsp.UpdateInfo = &userInfoProto.UpdateInfo{}
	DB := db.MasterEngine()
	err := DB.Model(&userInfoProto.UpdateInfo{}).First(&rsp.UpdateInfo).Error
	if err != nil {
		rsp.Error = &userInfoProto.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		return
	}
	rsp.Error = &userInfoProto.Error{Code: http.StatusOK, Message: "ok"}
	return
}
