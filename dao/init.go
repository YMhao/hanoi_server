package dao

import (
	"github.com/mediocregopher/radix.v2/redis"
	mgo "gopkg.in/mgo.v2"
)

var (
	dbname                          = "hanoi_api"
	UserDao        *_UserDao        = nil
	UserProfileDao *_UserProfileDao = nil
	UserBindDao    *_UserBindDao    = nil
	UserPasswdDao  *_UserPasswdDao  = nil
	SessionDao     *_SessionDao     = nil
)

func Init(mongoUrl, redisUrl string) error {
	sess, err := mgo.Dial(mongoUrl)
	if err != nil {
		return err
	}
	UserDao = newUserDao(sess)
	UserProfileDao = newUserProfileDao(sess)
	UserBindDao = newUserBindDao(sess)
	UserPasswdDao = newUserPasswdDao(sess)

	redisClient, err := redis.Dial("tcp", redisUrl)
	if err != nil {
		return err
	}
	SessionDao = newSessionDao(redisClient)
	return nil
}
