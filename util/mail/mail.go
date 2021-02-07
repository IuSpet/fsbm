package mail

import (
	"github.com/jordan-wright/email"
	"net/smtp"
)

var mailConfig *smtpConfig

type DefaultMail struct {
	Dest    []string
	Subject string
	Text    []byte
}

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

func newMail(h *DefaultMail) *email.Email {
	e := email.NewEmail()
	e.To = h.Dest
	e.Subject = h.Subject
	e.Text = h.Text
	return e
}

func SendMail(msg *DefaultMail) error {
	e := newMail(msg)
	e.From = mailConfig.from
	return e.Send(mailConfig.smtpAddr, smtp.PlainAuth("", mailConfig.username, mailConfig.passport, mailConfig.host))
}
