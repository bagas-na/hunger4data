package service

import (
	"context"
	"errors"
	"notification/internal/adapters/db"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockMailer to simulate the SMTP adapter
type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) Send(ctx context.Context, email Email) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func setupTestRepo(t *testing.T) (*db.NotificationLogRepo, *gorm.DB) {
	gormDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = gormDB.AutoMigrate(&db.NotificationLog{})
	require.NoError(t, err)

	return db.NewNotificationLogRepo(gormDB), gormDB
}

func TestEmailService_SendNotification(t *testing.T) {
	t.Run("Success - Email sends and logs as 'sent'", func(t *testing.T) {
		repo, gormDB := setupTestRepo(t)
		mockMailer := new(MockMailer)
		svc := NewNotificationService(mockMailer, repo)

		testEmail := Email{ToEmail: "user@test.com", Subject: "Hello"}
		mockMailer.On("Send", mock.Anything, testEmail).Return(nil)

		// Act
		err := svc.SendNotification(context.Background(), testEmail)

		// Assert
		assert.NoError(t, err)

		var savedLog db.NotificationLog
		gormDB.First(&savedLog)
		assert.Equal(t, "sent", savedLog.Status)
		assert.Equal(t, "user@test.com", savedLog.ToEmail)
		assert.Nil(t, savedLog.Error)
	})

	t.Run("Failure - Mailer fails and logs as 'failed' with error message", func(t *testing.T) {
		repo, gormDB := setupTestRepo(t)
		mockMailer := new(MockMailer)
		svc := NewNotificationService(mockMailer, repo)

		smtpError := errors.New("connection timeout")
		testEmail := Email{ToEmail: "fail@test.com", Subject: "Urgent"}
		mockMailer.On("Send", mock.Anything, testEmail).Return(smtpError)

		// Act
		err := svc.SendNotification(context.Background(), testEmail)

		// Assert
		assert.NoError(t, err)

		var savedLog db.NotificationLog
		gormDB.First(&savedLog)
		assert.Equal(t, "failed", savedLog.Status)
		assert.Equal(t, "connection timeout", *savedLog.Error)
	})
}
