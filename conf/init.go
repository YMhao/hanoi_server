package conf

import (
	"errors"

	"github.com/astaxie/beego/config"
)

type Configuration struct {
	MongoURL   string
	RedisURL   string
	Debug      bool
	ListenAddr string
}

func NewConfig(fileName string) (*Configuration, error) {
	configer, err := config.NewConfig("ini", fileName)
	if err != nil {
		return nil, err
	}
	conf := &Configuration{}
	conf.MongoURL = configer.String("mongodb.url")
	if conf.MongoURL == "" {
		return nil, errors.New("mongodb.url is missing")
	}
	conf.RedisURL = configer.String("redis.url")
	if conf.RedisURL == "" {
		return nil, errors.New("redis.url is missing")
	}
	conf.Debug = configer.DefaultBool("debug", false)
	conf.ListenAddr = configer.String("listen.addr")
	if conf.ListenAddr == "" {
		return nil, errors.New("listen.addr is missing")
	}
	return conf, nil
}
