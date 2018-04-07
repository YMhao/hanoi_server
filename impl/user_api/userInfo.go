package user_api

import (
	"github.com/YMhao/EasyApi/serv"
	"github.com/YMhao/hanoi_server/dao"
	"github.com/gin-gonic/gin"
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
	SetUserInfoBack,
)

func SetUserInfoBack(data []byte, c *gin.Context) (interface{}, *serv.APIError) {
	req := &SetUserInfoRequest{}
	err := serv.UnmarshalAndCheckValue(data, req)
	if err != nil {
		return nil, serv.NewError(err)
	}
	userUUID, err := dao.SessionDao.GetUserUUID(req.SessionID)
	if err != nil {
		return nil, serv.NewError(err)
	}
	dao.UserProfileDao.Set(&dao.UserProfile{
		UserUUID: userUUID,
		NickName: req.UserInfo.NickName,
		Gender:   req.UserInfo.Gender,
		//BirthDayTimeStamp : req.UserInfo.BirthDay,
		MobilePhone: req.UserInfo.Mobilephone,
		Email:       req.UserInfo.Email,
	})
	return &SetUserInfoRespone{}, nil
}
