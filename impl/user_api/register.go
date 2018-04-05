package user_api

import (
	"github.com/YMhao/EasyApi/serv"
)

type RegisterUserRequest struct {
	Type     string `desc:"类型" enum:"EMAL,MOBLEPHONE"`
	UserName string `desc:"用户名"`
	Passwd   string `desc:"密码，md5(passwd)"`
}

type RegisterUserResponse struct {
	Message string `desc:"提示信息，例如“注册成功，已向邮箱发送绑定帐号信息邮件，点击邮箱中的链接即可完成绑定”"`
}

var RegisterUserApi = serv.NewAPI(
	"registerUser",
	"注册新用户接口",
	&RegisterUserRequest{},
	&RegisterUserResponse{},
	nil,
)
