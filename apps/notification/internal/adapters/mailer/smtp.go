package mailer

import (
	"context"
	"fmt"
	"notification/internal/service"

	"gopkg.in/gomail.v2"
)

type SMTPMailer struct {
	host      string
	port      int
	username  string
	password  string
	fromName  string
	fromEmail string
}

func NewSMTPMailer(host string, port int, username, password, fromName, fromEmail string) *SMTPMailer {
	return &SMTPMailer{
		host:      host,
		port:      port,
		username:  username,
		password:  password,
		fromName:  fromName,
		fromEmail: fromEmail,
	}
}

func (s *SMTPMailer) Send(ctx context.Context, email service.Email) error {
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", s.fromEmail, s.fromName)
	msg.SetAddressHeader("To", email.ToEmail, email.ToName)
	msg.SetHeader("Subject", email.Subject)
	msg.SetBody("text/html", email.Body)

	d := gomail.NewDialer(s.host, s.port, s.username, s.password)

	// gomail does not support context directly, so we respect cancellation manually
	done := make(chan error, 1)
	go func() {
		err := d.DialAndSend(msg)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		if err != nil {
			return fmt.Errorf("smtp send failed: %w", err)
		}
	}

	return nil
}
