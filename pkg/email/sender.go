package email

import (
	"errors"
)

type Send struct {
	Recipient string
	Subject   string
	Body      string
}

type Sender interface {
	Send(input Send) error
}

func (e *Send) Validate() error {
	if e.Recipient == "" {
		return errors.New("empty to")
	}

	if e.Subject == "" || e.Body == "" {
		return errors.New("empty subject/body")
	}

	if !IsEmailValid(e.Recipient) {
		return errors.New("invalid to email")
	}

	return nil
}
