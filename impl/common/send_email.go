package common

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/YMhao/hanoi_server/conf"
)

var EmailSetting *conf.EmailSetting = nil

func InitEmailSetting(setting *conf.EmailSetting) {
	EmailSetting = setting
}

func SendSignUpSuccessEmail(toUser string) error {
	auth := smtp.PlainAuth("", EmailSetting.UserName, EmailSetting.Passwd, EmailSetting.Host)
	to := []string{toUser}
	nickname := EmailSetting.NickName
	user := EmailSetting.UserName
	subject := "Sign up Success! 注册成功"
	content_type := "Content-Type: text/plain; charset=UTF-8"
	body := fmt.Sprintf("注册成功，祝您游戏愉快。温馨提示：帐号已绑定您的邮箱，如忘记密码，可通过该邮箱找回密码。")
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	return smtp.SendMail(EmailSetting.Addr, auth, user, to, msg)
}
