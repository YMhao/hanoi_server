package user_api

import (
	"errors"
	"fmt"
	"regexp"

	"gopkg.in/mgo.v2"

	"github.com/YMhao/EasyApi/serv"
	"github.com/YMhao/hanoi_server/dao"
	"github.com/YMhao/hanoi_server/impl/common"
	"github.com/gin-gonic/gin"
)

type SignInOrSignUpRequest struct {
	UserName string `desc:"用户名"`
	Passwd   string `desc:"密码，md5(passwd), 32位md5"`
	IMIE     string `desc:"IMIE识别码,限制疯狂注册"`
}

type SignInOrSignUpResponse struct {
	HasLock     bool         `desc:"帐号是否被锁定"`
	LockTime    int64        `desc:"锁定到什么时候, 时间戳，单位是秒"`
	TryCount    int          `desc:"剩余重试次数"`
	PasswdError bool         `desc:"密码是否错误"`
	SessionID   string       `desc:"会话id"`
	SendEmail   bool         `desc:"是否发送邮件成功"`
	SendMessage bool         `desc:"是否发送短信成功"`
	NewUser     bool         `desc:"是否是新用户"`
	UserType    UserNameType `desc:"用户类型：邮箱、移动电话、未知" enum:"EMAIL、MOBILE_PHTONE、UNKNOWN"`
}

var SignInOrSignUpApi = serv.NewAPI(
	"SignInOrSignUp",
	`登录或注册用户接口`,
	&SignInOrSignUpRequest{},
	&SignInOrSignUpResponse{},
	SignInOrSignUpCallBack,
)

func SignInOrSignUpCallBack(data []byte, c *gin.Context) (interface{}, *serv.APIError) {
	req := &SignInOrSignUpRequest{}
	err := serv.UnmarshalAndCheckValue(data, req)
	if err != nil {
		return nil, serv.NewError(err)
	}

	userNameType := getUserNameType(req.UserName)
	if userNameType == USER_NAME_TYPE_UNKNOWN {
		return nil, &serv.APIError{
			Code:    "format.invalid",
			Message: "用户名格式错误",
		}
	}

	ok := checkPasswdFormat(req.Passwd)
	if !ok {
		return nil, &serv.APIError{
			Code:    "format.invalid",
			Message: "密码格式错误",
		}
	}

	userUUID, err := dao.UserDao.GetUUIDByName(req.UserName)
	if err != nil {
		if err != mgo.ErrNotFound {
			return nil, serv.NewError(err)
		}
		return signUp(req.UserName, req.Passwd, userNameType)
	}
	return signIn(userUUID, req.Passwd, userNameType)
}

func signIn(userUUID, passwd string, userNameType UserNameType) (interface{}, *serv.APIError) {
	ok, lockTime, tryCount, err := dao.UserPasswdDao.Check(userUUID, passwd)
	if err != nil {
		return nil, serv.NewError(err)
	}
	if !ok {
		return &SignInOrSignUpResponse{
				HasLock:     (lockTime > 0),
				LockTime:    lockTime,
				TryCount:    tryCount,
				SessionID:   "",
				PasswdError: true,
			}, &serv.APIError{
				Code:    "wrong.passwd",
				Message: "密码错误",
			}
	}

	sessionID, err := dao.SessionDao.NewSession(userUUID)
	if err != nil {
		return nil, serv.NewError(err)
	}
	return &SignInOrSignUpResponse{
		SessionID: sessionID,
		TryCount:  tryCount,
	}, nil
}

func getUserNameType(userName string) UserNameType {
	reg := `^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	rgx := regexp.MustCompile(reg)
	ok := rgx.MatchString(userName)
	if ok {
		return USER_NAME_TYPE_EMAIL
	}

	reg = `^1([38][0-9]|14[57]|5[^4])\d{8}$`
	rgx = regexp.MustCompile(reg)
	ok = rgx.MatchString(userName)
	if ok {
		return USER_NAME_TYPE_MOBILE_PHONE
	}
	return USER_NAME_TYPE_UNKNOWN
}

func checkPasswdFormat(passwd string) bool {
	reg := "^[0-9a-zA-Z]{32}$"
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(passwd)
}

func signUp(userName string, passwd string, userNameType UserNameType) (interface{}, *serv.APIError) {
	userUUID, err := dao.UserDao.Create(userName)
	if err != nil {
		return nil, serv.NewError(err)
	}

	err = dao.UserPasswdDao.Create(userUUID, passwd)
	if err != nil {
		return nil, serv.NewError(err)
	}

	respone := &SignInOrSignUpResponse{
		TryCount: dao.UserPasswdDao.TryCountMax(),
		NewUser:  true,
		UserType: userNameType,
	}
	// tryCreatBind(userName, userName, userNameType)
	err = trySendBindMessage(userName, userNameType)
	setSendStatus(respone, userNameType, err)
	sessionID, err := dao.SessionDao.NewSession(userUUID)
	if err != nil {
		return nil, serv.NewError(err)
	}
	respone.SessionID = sessionID
	return respone, nil
}

func setSendStatus(response *SignInOrSignUpResponse, userNameType UserNameType, err error) {
	if err != nil {
		return
	}
	switch userNameType {
	case USER_NAME_TYPE_EMAIL:
		response.SendEmail = true
	case USER_NAME_TYPE_MOBILE_PHONE:
		// TODO in future
	}
}

func tryCreatBind(userName, content string, userNameType UserNameType) error {
	userUUID, err := dao.UserDao.GetUUIDByName(userName)
	if err != nil {
		return err
	}

	bindInfo := &dao.UserBind{
		UserUUID: userUUID,
	}

	switch userNameType {
	case USER_NAME_TYPE_EMAIL:
		bindInfo.Email = content
	case USER_NAME_TYPE_MOBILE_PHONE:
		bindInfo.MobilePhone = content
	}

	err = dao.UserBindDao.Create(bindInfo)
	if err != nil {
		return err
	}
	return nil
}

func trySendBindMessage(dest string, userNameType UserNameType) error {
	switch userNameType {
	case "EMAIL":
		err := common.SendSignUpSuccessEmail(dest)
		if err != nil {
			fmt.Println("send Email", err)
		}
		return err
	case "MOBILEPHONE":
		// TODO in future
		return errors.New("SMS are not currently supported")
	}
	return nil
}
