package user_info

import (
	"book-user_api/m_client"
	"context"
	"errors"
	"gitee.com/qianxunke/book-ticket-common/basic"
	"gitee.com/qianxunke/book-ticket-common/basic/api_common"
	"gitee.com/qianxunke/book-ticket-common/basic/common"
	"gitee.com/qianxunke/book-ticket-common/proto/auth"
	"gitee.com/qianxunke/book-ticket-common/proto/user"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"log"
	"net/http"
)

func Init(client client.Client) *UserApiService {
	return &UserApiService{
		serviceClient: user.NewUserInfoService(basic.UserService, client),
	}
}

type UserApiService struct {
	serviceClient user.UserInfoService
}

//登陆
func (userApiService *UserApiService) Login(c *gin.Context) {
	var reqInLogin user.InDoneUserLogin
	if err := c.ShouldBindJSON(&reqInLogin); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("[Api] 请求参数不合法！"))
		return
	}
	log.Printf("reqInLogin : %v\n", reqInLogin)
	//调用后台服务
	rsp, _ := userApiService.serviceClient.DoneUserLogin(context.TODO(), &reqInLogin)
	//返回结果
	response := &api_common.ResponseEntity{}
	if rsp.Error.Code == http.StatusOK {
		//将token写到cookies中去
		//	c.Writer.Header().Add("Content-Type", "application/json;charset=utf-8")
		c.Writer.Header().Add(common.RememberMeCookieName, rsp.Token)
		// 过期30分钟
		c.SetCookie(common.RememberMeCookieName, rsp.Token, 90000, "/", "", false, false)
		data := map[string]interface{}{}
		data["token"] = rsp.Token
		data["user"] = rsp.UserInf
		response.Message = rsp.Error.Message
		response.Code = http.StatusOK
		response.Data = data
		c.AbortWithStatusJSON(http.StatusOK, response)
	} else {
		response.Message = rsp.Error.Message
		response.Code = http.StatusBadRequest
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
}

type Register struct {
	Nike_name         string `json:"nike_name"`
	User_name         string `json:"user_name"`
	Password          string `json:"password"`
	Mobile_phone      string `json:"mobile_phone"`
	User_email        string `json:"user_email"`
	Verification_code string `json:"verification_code"`
}

//注册
func (userApiService *UserApiService) Register(c *gin.Context) {

	register := &Register{}
	if err := c.ShouldBindJSON(&register); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("[Api] 请求参数不合法！"))
		return
	}
	reqInRegister := user.InDoneUserRegister{Userinf: &user.UserInf{UserName: register.User_name, UserEmail: register.User_email, MobilePhone: register.Mobile_phone, Password: register.Password},
		VerificationCode: register.Verification_code}
	//返回结果
	response := &api_common.ResponseEntity{}
	if reqInRegister.Userinf == nil || len(reqInRegister.VerificationCode) == 0 {
		api_common.SrvResultDone(c, nil, &api_common.Error{Code: http.StatusBadRequest, Message: "请求体为空"})
		return
	}

	//判断基本信息是否合法
	if len(reqInRegister.Userinf.MobilePhone) == 0 {
		api_common.SrvResultDone(c, nil, &api_common.Error{Code: http.StatusBadRequest, Message: "电话号码不能为空"})
		return
	}

	if len(reqInRegister.VerificationCode) == 0 {
		api_common.SrvResultDone(c, nil, &api_common.Error{Code: http.StatusBadRequest, Message: "验证码不能为空"})
		return
	}

	if len(reqInRegister.Userinf.UserName) == 0 {
		api_common.SrvResultDone(c, nil, &api_common.Error{Code: http.StatusBadRequest, Message: "用户名不能为空"})
		return
	}
	if len(reqInRegister.Userinf.Password) == 0 {
		api_common.SrvResultDone(c, nil, &api_common.Error{Code: http.StatusBadRequest, Message: "密码不能为空"})
		return
	}
	//调用后台服务
	rsp, _ := userApiService.serviceClient.DoneUserRegister(context.TODO(), &user.InDoneUserRegister{
		VerificationCode: reqInRegister.VerificationCode,
		Userinf:          reqInRegister.Userinf,
	})
	if rsp.Error.Code == http.StatusOK {
		//将token写到cookies中去
		c.Writer.Header().Add("set-cookie", "application/json; charset=utf-8")
		// 过期30分钟
		c.SetCookie(common.RememberMeCookieName, rsp.Token, 90000, "/", "", false, false)
		data := map[string]interface{}{}
		data["token"] = rsp.Token
		data["user"] = rsp.UserInf
		response.Message = rsp.Error.Message
		response.Code = http.StatusOK
		response.Data = data
		c.Writer.Header().Add("Content-Type", "application/json; charset=utf-8")
		c.AbortWithStatusJSON(http.StatusOK, response)
	} else {
		api_common.SrvResultDone(c, nil, &api_common.Error{Code: rsp.Error.Code, Message: rsp.Error.Message})
	}
}

// 退出登录
func (userApiService *UserApiService) Logout(c *gin.Context) {
	response := &api_common.ResponseEntity{}
	token, _ := c.Cookie(common.RememberMeCookieName)
	if len(token) == 0 {
		response.Message = "token失效"
		response.Code = http.StatusBadRequest
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var err error
	_, err = m_client.AuthClient.DelUserAccessToken(context.TODO(), &auth.Request{
		Token: token,
	})
	if err != nil {
		response.Message = "退出登陆失败！"
		response.Code = http.StatusInternalServerError
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	// 清除cookie
	c.SetCookie(common.RememberMeCookieName, "", 0, "/", "", false, false)
	// 返回JSON结构
	response.Code = http.StatusOK
	response.Message = "退出登陆成功"
	c.AbortWithStatusJSON(http.StatusOK, response)
}

//获取验证码
func (userApiService *UserApiService) GetCode(c *gin.Context) {
	requestParams := &user.InGetVerificationCode{}
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("[Api] 请求参数不合法！"))
		return
	}
	if len(requestParams.Telephone) == 0 {
		api_common.SrvResultDone(c, nil, &api_common.Error{Code: http.StatusBadRequest, Message: "手机号码不能为空"})
		return
	}
	log.Printf("requestParams : %v\n", requestParams)
	rsp, _ := userApiService.serviceClient.GetVerificationCode(context.TODO(), requestParams)
	api_common.SrvResultDone(c, nil, &api_common.Error{Code: rsp.Error.Code, Message: rsp.Error.Message})
}

//获取用户列表
func (userApiService *UserApiService) GetUserInfoList(c *gin.Context) {
	//	requestParams := &user.InGetUserInfoList{}
	/*
		if err := c.ShouldBindJSON(&requestParams); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("[Api] 请求参数不合法！"))
			return
		}

	*/

	_ = c.AbortWithError(http.StatusBadRequest, errors.New("[Api] 请求参数不合法！"))

	//	rsp, _ := userApiService.serviceClient.GetUserInfoList(context.TODO(), requestParams)
	//	api_common.SrvResultListDone(c, rsp.UserInfList, rsp.Limit, rsp.Pages, rsp.Total, &api_common.Error{Code: rsp.Error.Code, Message: rsp.Error.Message})
}

//修改用户普通信息
func (userApiService *UserApiService) UpdateUserInfo(c *gin.Context) {
	requestParams := &user.InUpdateUserInfo{}
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("[Api] 请求参数不合法！"))
		return
	}
	var userId string
	if userId := api_common.GetHeadUserId(c); len(userId) == 0 {
		api_common.SrvResultDone(c, nil, &api_common.Error{Code: http.StatusBadRequest, Message: "身份过期，请重新登陆"})
		return
	}
	requestParams.UserInf.UserId = userId
	rsp, _ := userApiService.serviceClient.UpdateUserInfo(context.TODO(), requestParams)
	api_common.SrvResultDone(c, rsp.UserInf, &api_common.Error{Code: rsp.Error.Code, Message: rsp.Error.Message})
}

//获取用户信息
func (userApiService *UserApiService) GetUserInfo(c *gin.Context) {
	requestParams := &user.InGetUserInfo{}
	requestParams.UserId = api_common.GetHeadUserId(c)
	if len(requestParams.UserId) == 0 {
		api_common.SrvResultDone(c, nil, &api_common.Error{Code: http.StatusBadRequest, Message: "身份过期，请重新登陆"})
		return
	}
	rsp, _ := userApiService.serviceClient.GetUserInfo(context.TODO(), requestParams)

	api_common.SrvResultDone(c, rsp, &api_common.Error{Code: rsp.Error.Code, Message: rsp.Error.Message})
}

//修改用户普通信息
func (userApiService *UserApiService) Login12306(c *gin.Context) {
	requestParams := &user.In_Login12306{}
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("[Api] 请求参数不合法！"))
		return
	}
	requestParams.UserId = api_common.GetHeadUserId(c)
	if len(requestParams.UserId) == 0 {
		api_common.SrvResultDone(c, nil, &api_common.Error{Code: http.StatusBadRequest, Message: "身份过期，请重新登陆"})
		return
	}
	rsp, _ := userApiService.serviceClient.Login12306(context.TODO(), requestParams)
	api_common.SrvResultDone(c, nil, &api_common.Error{Code: rsp.Error.Code, Message: rsp.Error.Message})
}

//获取用户列表
func (userApiService *UserApiService) GetUserPresenters(c *gin.Context) {
	requestParams := &user.In_GetUserPassengerList{}
	userId := api_common.GetHeadUserId(c)
	if len(userId) == 0 {
		api_common.SrvResultDone(c, nil, &api_common.Error{Code: http.StatusBadRequest, Message: "身份过期，请重新登陆"})
		return
	}
	requestParams.UserId = userId
	rsp, _ := userApiService.serviceClient.GetUserPassengerList(context.TODO(), requestParams)
	api_common.SrvResultListDone(c, rsp.PassengerList, 0, 0, 0, &api_common.Error{Code: rsp.Error.Code, Message: rsp.Error.Message})
}

func (userApiService *UserApiService) AppUpdateInfo(c *gin.Context) {
	requestParams := &user.In_UpdateInfo{}

	rsp, _ := userApiService.serviceClient.GetUpdateInfo(context.TODO(), requestParams)
	api_common.SrvResultDone(c, rsp.UpdateInfo, &api_common.Error{Code: rsp.Error.Code, Message: rsp.Error.Message})
}
