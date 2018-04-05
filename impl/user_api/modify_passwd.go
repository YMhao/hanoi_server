package user_api

import (
	"github.com/YMhao/EasyApi/serv"
)

type ModifyPasswdRequest struct {
	UserName  string `desc:"用户名"`
	OldPasswd string `desc:"旧密码，md5(passwd)"`
	NewPasswd string `desc:"新密码，md5(passwd)"`
}

type ModifyPasswdResponse struct {
}

var ModifyPasswdApi = serv.NewAPI(
	"ModifyPasswd",
	"修改密码",
	&ModifyPasswdRequest{},
	&ModifyPasswdResponse{},
	nil,
)
