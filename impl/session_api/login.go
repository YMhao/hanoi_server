package session_api

import (
	"strings"

	"github.com/YMhao/EasyApi/serv"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type LoginRequest struct {
	UserName string `desc:"用户名"`
	Passwd   string `desc:"密码，md5加密"`
}

type LoginRespone struct {
	SessionID string `desc:"会话id"`
}

var LoginAPI = serv.NewAPI(
	"login",
	`登录接口`,
	&LoginRequest{},
	&LoginRespone{},
	HelloAPICall,
)

func newSession() (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return strings.Replace(uid.String(), "-", "", -1), nil
}

func HelloAPICall(data []byte, c *gin.Context) (interface{}, *serv.APIError) {
	req := &LoginRequest{}
	err := serv.UnmarshalAndCheckValue(data, req)
	if err != nil {
		return nil, serv.NewError(err)
	}

	sessionId, err := newSession()
	if err != nil {
		return nil, serv.NewError(err)
	}

	return &LoginRespone{
		SessionID: sessionId,
	}, nil
}
