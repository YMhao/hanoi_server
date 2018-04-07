package common

import (
	"fmt"

	"github.com/YMhao/hanoi_server/conf"
	gomail "gopkg.in/gomail.v2"
)

var EmailSetting *conf.EmailSetting = nil

func InitEmailSetting(setting *conf.EmailSetting) {
	EmailSetting = setting
}

var messge = `
<pre>
您好！

谢谢您注册为汉诺塔游戏的会员!
祝您游戏愉快!

登录帐号和密码如下：

邮箱地址: %s
密码: 还记得您注册时你所输入的密码吗?
点击以下链接进行邮箱验证:</pre>
<a href="%s">%s</a>
`

func createMessage(userName, url string) string {
	return fmt.Sprintf(messge, userName, url, url)
}

func SendSignUpSuccessEmail(toUser string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", EmailSetting.UserName, EmailSetting.NickName)
	m.SetHeader("To", m.FormatAddress(toUser, "user"))
	m.SetHeader("Subject", "注册成功!欢迎光临汉诺塔游戏")
	m.SetBody("text/html", createMessage(toUser, "#"))

	d := gomail.NewPlainDialer(EmailSetting.Host, EmailSetting.SMTPPort, EmailSetting.UserName, EmailSetting.Passwd)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
