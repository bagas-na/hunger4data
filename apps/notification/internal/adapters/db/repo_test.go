package db

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn: sqlDB,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return db, mock
}

func TestNotificationLogRepo_Create(t *testing.T) {
	t.Run("success - should insert notification log", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewNotificationLogRepo(db)

		logID := uuid.New()
		logEntry := &NotificationLog{
			ID:      logID,
			ToEmail: "test@example.com",
			Subject: "Hello",
			Status:  "sent",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "notification_logs"`)).
			WithArgs(
				logEntry.ToEmail,
				logEntry.Subject,
				logEntry.Status,
				nil,              // Error
				sqlmock.AnyArg(), // CreatedAt
				logEntry.ID,      // ID
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(logID))

		mock.ExpectCommit()

		err := repo.Create(context.Background(), logEntry)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure - database constraint error", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewNotificationLogRepo(db)

		logEntry := &NotificationLog{ToEmail: "fail@example.com"}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "notification_logs"`)).
			WillReturnError(errors.New("null value in column \"status\" violates not-null constraint"))
		mock.ExpectRollback()

		err := repo.Create(context.Background(), logEntry)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
