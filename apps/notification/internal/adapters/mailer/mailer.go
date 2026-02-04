package mailer

import (
	"context"
	"fmt"
	"notification/internal/service"

	"github.com/mailersend/mailersend-go"
)

type MailersendMailer struct {
	client *mailersend.Mailersend
}

func NewMailersendMailer(apiKey string) *MailersendMailer {
	return &MailersendMailer{
		client: mailersend.NewMailersend(apiKey),
	}
}

func (m *MailersendMailer) Send(ctx context.Context, email service.Email) error {
	from := mailersend.From{
		Name:  email.FromName,
		Email: email.FromEmail,
	}

	recipient := mailersend.Recipient{
		Name:  email.ToName,
		Email: email.ToEmail,
	}

	msg := &mailersend.Message{
		From:       from,
		Recipients: []mailersend.Recipient{recipient},
		Subject:    email.Subject,
		HTML:       email.Body,
	}

	_, err := m.client.Email.Send(ctx, msg)
	if err != nil {
		return fmt.Errorf("mailersend send failed: %w", err)
	}

	return nil
}
