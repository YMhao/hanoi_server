package dao

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/satori/go.uuid"
	mgo "gopkg.in/mgo.v2"
)

type User struct {
	UserName string `bson:"_id"`
	UserUUID string `bson:"user_uuid"`
	Common   Common `bson:"common"`
}

type UserProfile struct {
	UserUUID          string `bson:"_id"`
	NickName          string `bson:"nick"`
	Gender            string `bson:"gender" enum:"MALE,FEMALE，UNKNOWN"`
	MobilePhone       string `bson:"mobile_phone"`
	Email             string `bson:"email"`
	BirthDayTimeStamp int64  `desc:"birthday"`
}

type UserBind struct {
	UserUUID    string `bson:"_id"`
	Email       string `bson:"email"`
	MobilePhone string `bson:"mobile_phone"`
}

type UserPasswd struct {
	UserUUID string `bson:"_id"`
	Passwd   string `bson:"passwd"` // md5(md5(passwd) + salt )
	Salt     string `bson:"salt"`
	LockTime int64  `bson:"lock_time"`
	TryCount int    `bson:"try_count"`
}

type _UserDao struct {
	sess *mgo.Session
	coll *mgo.Collection
}

func newUserDao(sess *mgo.Session) *_UserDao {
	return &_UserDao{
		sess: sess,
		coll: sess.DB(dbname).C("user"),
	}
}

func (d *_UserDao) Create(UserName string) (string, error) {
	c := d.coll.With(d.sess)
	t := time.Now()
	UUID, err := d.newUUID(t)
	if err != nil {
		return "", err
	}
	timeStamp := t.UnixNano() / 1e6
	return UUID, c.Insert(&User{
		UserName: UserName,
		UserUUID: UUID,
		Common: Common{
			Create: timeStamp,
			Update: timeStamp,
			Valid:  true,
		},
	})
}

func (d *_UserDao) IsUserExist(UserName string) (bool, error) {
	c := d.coll.With(d.sess)
	obj := &User{}
	err := c.FindId(UserName).One(obj)
	if err != nil {
		if err == mgo.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (d *_UserDao) GetUUIDByName(UserName string) (string, error) {
	c := d.coll.With(d.sess)
	obj := &User{}
	err := c.FindId(UserName).One(obj)
	if err != nil {
		return "", err
	}
	return obj.UserUUID, nil
}

func (d *_UserDao) newUUID(t time.Time) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	UUID := t.Format("20060102_150405_") + strings.Replace(id.String(), "-", "", -1)
	return UUID, nil
}

type _UserProfileDao struct {
	sess *mgo.Session
	coll *mgo.Collection
}

func newUserProfileDao(sess *mgo.Session) *_UserProfileDao {
	return &_UserProfileDao{
		sess: sess,
		coll: sess.DB(dbname).C("user_profile"),
	}
}

func (d *_UserProfileDao) Set(userProfile *UserProfile) error {
	c := d.coll.With(d.sess)
	_, err := c.UpsertId(userProfile.UserUUID, userProfile)
	return err
}

type _UserBindDao struct {
	sess *mgo.Session
	coll *mgo.Collection
}

func newUserBindDao(sess *mgo.Session) *_UserBindDao {
	return &_UserBindDao{
		sess: sess,
		coll: sess.DB(dbname).C("user_bind"),
	}
}

func (d *_UserBindDao) Create(info *UserBind) error {
	c := d.coll.With(d.sess)
	return c.Insert(info)
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func newSalt() string {
	salt := r.Intn(10000)
	return fmt.Sprintf("%d", salt)
}

func newPasswd(passwd string) (salt string, newPasswd string) {
	salt = newSalt()
	data := []byte(passwd + salt)
	has := md5.Sum(data)
	newPasswd = fmt.Sprintf("%x", has)
	return
}

func getPasswd(passwd, salt string) string {
	data := []byte(passwd + salt)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

type _UserPasswdDao struct {
	sess *mgo.Session
	coll *mgo.Collection
}

func newUserPasswdDao(sess *mgo.Session) *_UserPasswdDao {
	return &_UserPasswdDao{
		sess: sess,
		coll: sess.DB(dbname).C("user_passwd"),
	}
}

func (d *_UserPasswdDao) TryCountMax() int {
	return 8
}

func (d *_UserPasswdDao) Create(userUUID, passwd string) error {
	c := d.coll.With(d.sess)
	salt, newPasswd := newPasswd(passwd)
	return c.Insert(&UserPasswd{
		UserUUID: userUUID,
		Passwd:   newPasswd,
		Salt:     salt,
		LockTime: 0,
		TryCount: d.TryCountMax(),
	})
}

func (d *_UserPasswdDao) Update(userUUID, passwd string) error {
	c := d.coll.With(d.sess)
	obj, err := d.Get(userUUID)
	if err != nil {
		return err
	}
	salt, newPasswd := newPasswd(passwd)
	obj.Passwd = newPasswd
	obj.Salt = salt
	obj.LockTime = 0
	obj.TryCount = d.TryCountMax()
	return c.UpdateId(userUUID, obj)
}

func (d *_UserPasswdDao) CheckPasswOlny(userUUID, passwd string) (bool, error) {
	obj, err := d.Get(userUUID)
	if err != nil {
		return false, err
	}
	newPasswd := getPasswd(passwd, obj.Salt)
	return newPasswd == obj.Passwd, nil
}

func (d *_UserPasswdDao) Check(userUUID, passwd string) (bool, int64, int, error) {
	obj, err := d.Get(userUUID)
	if err != nil {
		return false, 0, 0, err
	}

	if obj.TryCount <= 0 {
		t := time.Now().UnixNano() / 1e6
		if t < obj.LockTime {
			return false, obj.LockTime, 0, nil
		}
		err = d.Recover(obj)
		if err != nil {
			return false, 0, 0, err
		}
	}

	newPasswd := getPasswd(passwd, obj.Salt)

	if newPasswd != obj.Passwd {
		d.Inc(obj)
		newObj, err := d.Get(obj.UserUUID)
		if err != nil {
			return false, 0, 0, err
		}
		return false, newObj.LockTime, newObj.TryCount, nil
	}

	d.Recover(obj)
	return true, 0, d.TryCountMax(), nil
}

func (d *_UserPasswdDao) Get(userUUID string) (*UserPasswd, error) {
	c := d.coll.With(d.sess)
	obj := &UserPasswd{}
	err := c.FindId(userUUID).One(obj)
	return obj, err
}

func (d *_UserPasswdDao) Recover(obj *UserPasswd) error {
	c := d.coll.With(d.sess)
	err := c.Update(bson.M{"_id": obj.UserUUID}, bson.M{
		"$set": bson.M{
			"lock_time": 0,
			"try_count": 5,
		},
	})
	if err != nil {
		return err
	}
	obj.LockTime = 0
	obj.TryCount = 10
	return nil
}

func (d *_UserPasswdDao) Inc(obj *UserPasswd) {
	c := d.coll.With(d.sess)
	if obj.TryCount <= 1 {
		err := c.Update(bson.M{"_id": obj.UserUUID}, bson.M{
			"$inc": bson.M{
				"try_count": -1,
			},
		})
		if err != nil {
			return
		}
		c.Update(bson.M{
			"_id":       obj.UserUUID,
			"lock_time": 0,
		}, bson.M{
			"$set": bson.M{
				"lock_time": time.Now().UnixNano()/1e6 + 1000*60*60*2,
			},
		})
	}
	return
}
