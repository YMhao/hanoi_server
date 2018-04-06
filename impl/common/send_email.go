package common

import (
	"github.com/YMhao/hanoi_server/conf"
	gomail "gopkg.in/gomail.v2"
)

var EmailSetting *conf.EmailSetting = nil

func InitEmailSetting(setting *conf.EmailSetting) {
	EmailSetting = setting
}

func SendSignUpSuccessEmail(toUser string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", EmailSetting.UserName, EmailSetting.NickName)
	m.SetHeader("To", m.FormatAddress(toUser, "user"))
	m.SetHeader("Subject", "汉诺塔")
	m.SetBody("text/html", "祝您游戏愉快！")

	d := gomail.NewPlainDialer(EmailSetting.Host, EmailSetting.SMTPPort, EmailSetting.UserName, EmailSetting.Passwd)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
