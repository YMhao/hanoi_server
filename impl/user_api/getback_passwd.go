package user_api

import (
	"errors"

	"github.com/YMhao/EasyApi/serv"
	"github.com/gin-gonic/gin"
)

type GetBackPasswdRequest struct {
	UserName string `desc:"用户名"`
	Type     string `desc:"找回密码的方式" enum:"EMAIL"`
}

type GetBackPasswdResponse struct {
	Type       string `desc:"找回密码的方式" enum:"EMAIL"`
	SendStatus bool   `desc:"是否已发送"`
}

var GetBackPasswdApi = serv.NewAPI(
	"GetBackPasswd",
	"找回密码",
	&GetBackPasswdRequest{},
	&GetBackPasswdResponse{},
	GetBackPasswdCallBack,
)

func GetBackPasswdCallBack(data []byte, c *gin.Context) (interface{}, *serv.APIError) {
	err := errors.New("暂不支持")
	return nil, serv.NewError(err)
}
