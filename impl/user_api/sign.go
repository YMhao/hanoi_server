package user_api

import (
	"regexp"

	"gopkg.in/mgo.v2"

	"github.com/YMhao/EasyApi/serv"
	"github.com/YMhao/hanoi_server/dao"
	"github.com/gin-gonic/gin"
)

type SignInOrSignUpRequest struct {
	UserName string `desc:"用户名"`
	Passwd   string `desc:"密码，md5(passwd), 32位md5"`
}

type SignInOrSignUpResponse struct {
	HasLock     bool   `desc:"帐号是否被锁定"`
	LockTime    int64  `desc:"锁定到什么时候, 时间戳，单位是毫秒"`
	TryCount    int    `desc:"剩余重试次数"`
	PasswdError bool   `desc:"密码是否错误"`
	SessionID   string `desc:"会话id"`
	Message     string `desc:"提示信息，例如“注册成功，已向邮箱发送绑定帐号信息邮件，点击邮箱中的链接即可完成绑定”"`
}

var SignInOrSignUpApi = serv.NewAPI(
	"SignInOrSignUp",
	"登录或注册用户",
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
		userUUID, err := dao.UserDao.Create(req.UserName)
		if err != nil {
			return nil, serv.NewError(err)
		}

		err = dao.UserPasswdDao.Create(userUUID, req.Passwd)
		if err != nil {
			return nil, serv.NewError(err)
		}

		respone := &SignInOrSignUpResponse{}
		err = tryCreatBind(req.UserName, req.UserName, userNameType)
		if err != nil {
			respone.Message = "未绑定"
		}
		trySendBindMessage(req.UserName, userNameType)

		sessionID, err := dao.SessionDao.NewSession(userUUID)
		if err != nil {
			return nil, serv.NewError(err)
		}
		return &SignInOrSignUpResponse{
			SessionID: sessionID,
		}, nil
	}

	ok, lockTime, tryCount, err := dao.UserPasswdDao.Check(userUUID, req.Passwd)
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
		}, nil
	}

	sessionID, err := dao.SessionDao.NewSession(userUUID)
	if err != nil {
		return nil, serv.NewError(err)
	}
	return &SignInOrSignUpResponse{
		SessionID: sessionID,
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

func trySendBindMessage(dest string, userNameType UserNameType) {
	switch userNameType {
	case "EMAIL":
		// TODO
	case "MOBILEPHONE":
		// TODO
	}
}
