package mail

import (
	"github.com/jordan-wright/email"
	"net/smtp"
)

var mailConfig *smtpConfig

type smtpConfig struct {
	from     string
	username string
	passport string
	smtpAddr string
	host     string
}

func init() {
	mailConfig = &smtpConfig{
		from:     "fsbm",
		username: "luSpet@163.com",
		passport: "CELSWCJKLYLATBOC",
		smtpAddr: "smtp.163.com:25",
		host:     "smtp.163.com",
	}
}

func newMail(dest []string, subject string, text []byte) *email.Email {
	e := email.NewEmail()
	e.To = dest
	e.Subject = subject
	e.Text = text
	return e
}

func SendMail(dest []string, subject string, text []byte) error {
	e := newMail(dest, subject, text)
	e.From = mailConfig.from
	return e.Send(mailConfig.smtpAddr, smtp.PlainAuth("", mailConfig.username, mailConfig.passport, mailConfig.host))
}
