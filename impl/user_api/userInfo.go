package user_api

import (
	"github.com/YMhao/EasyApi/serv"
)

type UserInfo struct {
	NickName    string `desc:"昵称"`
	Mobilephone string `desc:"手机号码"`
	Email       string `desc:"邮箱"`
	Gender      string `desc:"性别" enum:"MALE,FEMALE，UNKNOWN"`
	BirthDay    string `desc:"生日,格式2006.01.02"`
}

type SetUserInfoRequest struct {
	SessionID string   `desc:"会话id"`
	UserInfo  UserInfo `desc:"用户信息"`
}
type SetUserInfoRespone struct {
}

var SetUserInfoApi = serv.NewAPI(
	"SerUserInfo",
	"设置用户信息",
	&SetUserInfoRequest{},
	&SetUserInfoRespone{},
	nil,
)
