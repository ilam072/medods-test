package smtp

import (
	"errors"
	"gopkg.in/gomail.v2"
	"medods-test/pkg/email"
)

type SMTPSender struct {
	from string
	pass string
	host string
	port int
}

func NewSMTPSender(from, pass, host string, port int) (*SMTPSender, error) {
	if !email.IsEmailValid(from) {
		return nil, errors.New("invalid from email")
	}

	return &SMTPSender{from: from, pass: pass, host: host, port: port}, nil
}

func (s *SMTPSender) Send(input email.Send) error {
	if err := input.Validate(); err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", s.from)
	msg.SetHeader("Recipient", input.Recipient)
	msg.SetHeader("Subject", input.Subject)
	msg.SetBody("text/html", input.Body)

	dialer := gomail.NewDialer(s.host, s.port, s.from, s.pass)
	if err := dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
