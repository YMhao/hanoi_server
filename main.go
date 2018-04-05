package main

import (
	"github.com/YMhao/EasyApi/serv"
	"github.com/YMhao/hanoi_server/impl/security_api"
	"github.com/YMhao/hanoi_server/impl/session_api"
	"github.com/YMhao/hanoi_server/impl/user_api"
)

func main() {
	conf := serv.NewAPIServConf("1.0", "", "hannoi", "hannoi 服务")
	conf.DebugOn = true
	conf.ListenAddr = ":8089"

	setsOfAPIs := serv.APISets{
		"session": []serv.API{
			session_api.LoginAPI,
		},
		"user": []serv.API{
			user_api.RegisterUserApi,
			user_api.SetUserInfoApi,
			user_api.ModifyPasswdApi,
			user_api.GetBackPasswdApi,
		},
		"security": []serv.API{
			security_api.BindApi,
		},
	}

	serv.RunAPIServ(conf, setsOfAPIs)
}
