package service

import (
	"context"
	"fmt"
	"notification/internal/adapters/db"
	"time"

	"github.com/google/uuid"
)

type Email struct {
	FromName  string
	FromEmail string
	ToName    string
	ToEmail   string
	Subject   string
	Body      string
}

type Mailer interface {
	Send(ctx context.Context, email Email) error
}

type EmailService interface {
	SendNotification(ctx context.Context, email Email) error
}

type emailService struct {
	mailer Mailer
	repo   *db.NotificationLogRepo
}

func NewNotificationService(mailer Mailer, repo *db.NotificationLogRepo) EmailService {
	return &emailService{
		mailer: mailer,
		repo:   repo,
	}
}

func (s *emailService) SendNotification(ctx context.Context, email Email) error {
	log := db.NotificationLog{
		ID:        uuid.New(),
		ToEmail:   email.ToEmail,
		Subject:   email.Subject,
		Status:    "sent",
		Error:     nil,
		CreatedAt: time.Now().UTC(),
	}

	err := s.mailer.Send(ctx, email)

	if err != nil {
		errString := err.Error()

		log.Status = "failed"
		log.Error = &errString

		fmt.Printf("error sending notification: %v", err)

		return s.repo.Create(ctx, &log)
	}

	return s.repo.Create(ctx, &log)
}
