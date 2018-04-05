package user_api

import (
	"github.com/YMhao/EasyApi/serv"
)

type GetBackPasswdRequest struct {
	UserName string `desc:"用户名"`
	Type     string `desc:"找回密码的方式" enum:"EMAIL"`
}

type GetBackPasswdResponse struct {
	Message string `desc:"一些提示相关的信息"`
}

var GetBackPasswdApi = serv.NewAPI(
	"GetBackPasswd",
	"找回密码",
	&GetBackPasswdRequest{},
	&GetBackPasswdResponse{},
	nil,
)
