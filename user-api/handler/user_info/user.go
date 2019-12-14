package user_info

import (
	"context"
	"errors"
	"gitee.com/qianxunke/surprise-shop-common/basic"
	"gitee.com/qianxunke/surprise-shop-common/basic/api_common"
	"gitee.com/qianxunke/surprise-shop-common/basic/common"
	auth "gitee.com/qianxunke/surprise-shop-common/protos/auth"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"net/http"
	apiClient "surprise-shop-user_api/client"

	"gitee.com/qianxunke/surprise-shop-common/protos/user/user_info"
)

func Init(client client.Client) *UserApiService {
	return &UserApiService{
		serviceClient: user_info.NewUserInfoService(basic.UserService, client),
	}
}

type UserApiService struct {
	serviceClient user_info.UserInfoService
}

//登陆
func (userApiService *UserApiService) Login(c *gin.Context) {
	var reqInLogin user_info.InDoneUserLogin
	if err := c.ShouldBindJSON(&reqInLogin); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("[Api] 请求参数不合法！"))
		return
	}
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

//注册
func (userApiService *UserApiService) Register(c *gin.Context) {
	reqInRegister := &user_info.InDoneUserRegister{}
	if err := c.ShouldBindJSON(&reqInRegister); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("[Api] 请求参数不合法！"))
		return
	}
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
	rsp, _ := userApiService.serviceClient.DoneUserRegister(context.TODO(), &user_info.InDoneUserRegister{
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
	_, err = apiClient.AuthClient.DelUserAccessToken(context.TODO(), &auth.Request{
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
	requestParams := &user_info.InGetVerificationCode{}
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("[Api] 请求参数不合法！"))
		return
	}
	if len(requestParams.Telephone) == 0 {
		api_common.SrvResultDone(c, nil, &api_common.Error{Code: http.StatusBadRequest, Message: "手机号码不能为空"})
		return
	}
	rsp, _ := userApiService.serviceClient.GetVerificationCode(context.TODO(), requestParams)
	api_common.SrvResultDone(c, nil, &api_common.Error{Code: rsp.Error.Code, Message: rsp.Error.Message})
}

//获取用户列表
func (userApiService *UserApiService) GetUserInfoList(c *gin.Context) {
	requestParams := &user_info.InGetUserInfoList{}
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("[Api] 请求参数不合法！"))
		return
	}
	rsp, _ := userApiService.serviceClient.GetUserInfoList(context.TODO(), requestParams)
	api_common.SrvResultListDone(c, rsp.UserInfList, rsp.Limit, rsp.Pages, rsp.Total, &api_common.Error{Code: rsp.Error.Code, Message: rsp.Error.Message})
}

//修改用户普通信息
func (userApiService *UserApiService) UpdateUserInfo(c *gin.Context) {
	requestParams := &user_info.InUpdateUserInfo{}
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("[Api] 请求参数不合法！"))
		return
	}
	var userId int64
	if userId := api_common.GetHeadUserId(c); userId == -1 {
		return
	}
	requestParams.UserInf.UserId = userId
	rsp, _ := userApiService.serviceClient.UpdateUserInfo(context.TODO(), requestParams)
	api_common.SrvResultDone(c, rsp.UserInf, &api_common.Error{Code: rsp.Error.Code, Message: rsp.Error.Message})
}

//获取用户信息
func (userApiService *UserApiService) GetUserInfo(c *gin.Context) {
	requestParams := &user_info.InGetUserInfo{}
	if requestParams.UserId = api_common.GetHeadUserId(c); requestParams.UserId == -1 {
		api_common.SrvResultDone(c, nil, &api_common.Error{Code: http.StatusBadRequest, Message: "身份过期，请重新登陆"})
		return
	}
	rsp, _ := userApiService.serviceClient.GetUserInfo(context.TODO(), requestParams)

	api_common.SrvResultDone(c, rsp, &api_common.Error{Code: rsp.Error.Code, Message: rsp.Error.Message})
}
