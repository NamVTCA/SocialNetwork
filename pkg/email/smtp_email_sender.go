package email

import (
	"fmt"
	"net/smtp"
)

type SMTPEmailSender struct {
	host     string
	port     string
	username string
	password string
}

func NewSMTPEmailSender(host, port, username, password string) EmailSender {
	return &SMTPEmailSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (s *SMTPEmailSender) Send(to, subject, body string) error {
	from := s.username
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, to, subject, body)

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}
