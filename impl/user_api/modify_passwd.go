package user_api

import (
	"errors"

	"github.com/YMhao/EasyApi/serv"
	"github.com/YMhao/hanoi_server/dao"
	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
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
	ModifyPasswdCallBack,
)

func ModifyPasswdCallBack(data []byte, c *gin.Context) (interface{}, *serv.APIError) {
	req := &ModifyPasswdRequest{}
	err := serv.UnmarshalAndCheckValue(data, req)
	if err != nil {
		return nil, serv.NewError(err)
	}

	if !checkPasswdFormat(req.OldPasswd) || !checkPasswdFormat(req.NewPasswd) {
		return nil, &serv.APIError{
			Code:    "wrong.format.passwd",
			Message: "密码格式错误",
		}
	}

	userUUID, err := dao.UserDao.GetUUIDByName(req.UserName)
	if err != nil {
		if err == mgo.ErrNotFound {
			err = errors.New("用户不存在")
		}
		return nil, serv.NewError(err)
	}

	ok, err := dao.UserPasswdDao.CheckPasswOlny(userUUID, req.OldPasswd)
	if err != nil {
		return nil, serv.NewError(err)
	}
	if !ok {
		err = errors.New("密码错误")
		return nil, serv.NewError(err)
	}

	if req.OldPasswd == req.NewPasswd {
		err = errors.New("旧密码与新密码不能相同")
		return nil, serv.NewError(err)
	}
	err = dao.UserPasswdDao.Update(userUUID, req.NewPasswd)
	if err != nil {
		return nil, serv.NewError(err)
	}
	return &ModifyPasswdResponse{}, nil
}
