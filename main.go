package main

import (
	"fmt"
	"os"

	"github.com/YMhao/EasyApi/serv"
	"github.com/YMhao/hanoi_server/conf"
	"github.com/YMhao/hanoi_server/dao"
	"github.com/YMhao/hanoi_server/impl/security_api"
	"github.com/YMhao/hanoi_server/impl/user_api"
)

var (
	VERSION            = "v1"
	SERVER_NAME        = "hannoiAPIs"
	SERVER_DESCRIPTION = "service for a mini game - hannoi"
	BUILD_TIME         = ""
)

func help() {
	fmt.Println("Usage:", os.Args[0], "[configuration file]")
	os.Exit(1)
}

func exitIfError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) != 2 {
		help()
	}

	cfg, err := conf.NewConfig(os.Args[1])
	exitIfError(err)

	err = dao.Init(cfg.MongoURL, cfg.RedisURL)
	exitIfError(err)

	servCfg := serv.NewAPIServConf(VERSION, BUILD_TIME, SERVER_NAME, SERVER_DESCRIPTION)
	servCfg.DebugOn = cfg.Debug
	servCfg.ListenAddr = cfg.ListenAddr
	servCfg.HTTPProxy = cfg.HttpProxy

	setsOfAPIs := serv.APISets{
		"user": []serv.API{
			user_api.SignInOrSignUpApi,
			user_api.SetUserInfoApi,
			user_api.ModifyPasswdApi,
			user_api.GetBackPasswdApi,
		},
		"security": []serv.API{
			security_api.BindApi,
		},
	}

	serv.RunAPIServ(servCfg, setsOfAPIs)
}
