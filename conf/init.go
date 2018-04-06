package conf

import (
	"errors"

	"github.com/astaxie/beego/config"
)

type Configuration struct {
	MongoURL     string
	RedisURL     string
	Debug        bool
	ListenAddr   string
	HttpProxy    string
	EmailSetting EmailSetting
}

type EmailSetting struct {
	UserName string
	Passwd   string
	Host     string
	NickName string
	Addr     string
}

func newMissError(key string) error {
	return errors.New(key + " is missing")
}

func configString(cfg config.Configer, key string) (string, error) {
	value := cfg.String(key)
	if value == "" {
		return "", errors.New(key + " is missing")
	}
	return value, nil
}

func NewConfig(fileName string) (*Configuration, error) {
	configer, err := config.NewConfig("ini", fileName)
	if err != nil {
		return nil, err
	}
	conf := &Configuration{}
	conf.MongoURL, err = configString(configer, "mongodb.url")
	if err != nil {
		return nil, err
	}

	conf.RedisURL, err = configString(configer, "redis.url")
	if err != nil {
		return nil, err
	}
	conf.Debug = configer.DefaultBool("debug", false)

	conf.ListenAddr, err = configString(configer, "listen.addr")
	if err != nil {
		return nil, err
	}
	conf.HttpProxy = configer.String("http.proxy")

	conf.EmailSetting.UserName, err = configString(configer, "email.username")
	if err != nil {
		return nil, err
	}
	conf.EmailSetting.Passwd, err = configString(configer, "email.passwd")
	if err != nil {
		return nil, err
	}
	conf.EmailSetting.NickName, err = configString(configer, "email.nickname")
	if err != nil {
		return nil, err
	}
	conf.EmailSetting.Host, err = configString(configer, "email.host")
	if err != nil {
		return nil, err
	}
	conf.EmailSetting.Addr, err = configString(configer, "email.addr")
	if err != nil {
		return nil, err
	}
	return conf, nil
}
