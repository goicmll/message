package message

import (
	"strings"

	"gopkg.in/gomail.v2"
)

type mail struct {
	host     string
	port     int
	account  string
	password string
}

func NewMail(host string, port int, account, password string) *mail {
	return &mail{host, port, account, password}
}

func (m *mail) SendText(subject, from, toStr, CcStr, body string) error {
	d := gomail.NewDialer(m.host, m.port, m.account, m.password)
	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	msg := gomail.NewMessage(
		//发送文本时设置编码，防止乱码。 如果txt文本设置了之后还是乱码，那可以将原txt文本在保存时
		//就选择utf-8格式保存
		gomail.SetEncoding(gomail.Base64),
	)
	// 添加别名
	msg.SetHeader("From", msg.FormatAddress(m.account, from))
	// 发送给用户(可以多个)
	msg.SetHeader("To", strings.Split(toStr, ",")...)
	// 设置邮件主题
	msg.SetHeader("Subject", subject)
	msg.SetBody("text", body)
	msg.SetHeader("Cc", strings.Split(CcStr, ",")...)
	err := d.DialAndSend(msg)
	return err
}

func (m *mail) SendTextWithAttach(subject, from, toStr, CcStr, body string, filePath []string) error {

	d := gomail.NewDialer(m.host, m.port, m.account, m.password)
	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	msg := gomail.NewMessage(
		//发送文本时设置编码，防止乱码。 如果txt文本设置了之后还是乱码，那可以将原txt文本在保存时
		//就选择utf-8格式保存
		gomail.SetEncoding(gomail.Base64),
	)
	for _, fp := range filePath {
		msg.Attach(fp)
	}
	// 添加别名
	msg.SetHeader("From", msg.FormatAddress(m.account, from))
	// 发送给用户(可以多个)
	msg.SetHeader("To", strings.Split(toStr, ",")...)
	// 设置邮件主题
	msg.SetHeader("Subject", subject)
	msg.SetHeader("Cc", strings.Split(CcStr, ",")...)
	msg.SetBody("text", body)
	err := d.DialAndSend(msg)
	return err
}
