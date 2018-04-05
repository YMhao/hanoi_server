package security_api

import (
	"github.com/YMhao/EasyApi/serv"
)

type BindRequest struct {
	SessionID string `desc:"会话id"`
	UserName  string `desc:"用户名"`
	Passwd    string `desc:"密码，md5加密"`
	Type      string `desc:"绑定类型" enum:"BIND_EMAIL,UNBIND_EMAIL"`
	Content   string `desc:"绑定的内容"`
}

type BindResponse struct {
	Message string `desc:"一些提示相关的反馈信息"`
}

var BindApi = serv.NewAPI(
	"bindOrUnbind",
	"绑定或解绑定",
	&BindRequest{},
	&BindResponse{},
	nil,
)
