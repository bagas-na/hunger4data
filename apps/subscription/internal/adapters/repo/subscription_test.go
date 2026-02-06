package repo

import (
	"errors"
	"regexp"
	"subscription/internal/adapters/model"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	return gormDB, mock, err
}

func TestCreateSubscription(t *testing.T) {
	db, mock, _ := SetupMockDB()
	repo := NewSubRepo(db)
	t.Run("Successful creation", func(t *testing.T) {
		UserID := uuid.New()
		subs := model.Subscription{
			UserID:      UserID,
			CountryCode: "AFG",
		}
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "subscriptions"`).
			WithArgs(
				sqlmock.AnyArg(),
				subs.UserID,
				subs.CountryCode,
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				sqlmock.AnyArg(), // DeletedAt
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.CreateSubcription(subs)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Failed creation", func(t *testing.T) {
		UserID := uuid.New()
		subs := model.Subscription{
			UserID:      UserID,
			CountryCode: "AFG",
		}
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "subscriptions"`).
			WithArgs(
				sqlmock.AnyArg(),
				subs.UserID,
				subs.CountryCode,
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				sqlmock.AnyArg(), // DeletedAt
			).
			WillReturnError(errors.New("database connection lost"))

		mock.ExpectRollback()

		err := repo.CreateSubcription(subs)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection lost")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Failure - Duplicate Key (23505)", func(t *testing.T) {
		input := model.Subscription{
			UserID:      uuid.New(),
			CountryCode: "CAF",
		}
		pgErr := &pgconn.PgError{Code: "23505"}

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "subscriptions"`).
			WithArgs(sqlmock.AnyArg(), input.UserID, input.CountryCode, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(pgErr)
		mock.ExpectRollback()

		err := repo.CreateSubcription(input)
		assert.ErrorIs(t, err, gorm.ErrDuplicatedKey)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetSubscriptionsByUserID(t *testing.T) {
	db, mock, _ := SetupMockDB()
	repo := NewSubRepo(db)
	t.Run("Success - Multiple subscriptions found", func(t *testing.T) {
		userID := uuid.New()
		rows := sqlmock.NewRows([]string{"id", "user_id", "country_code", "created_at", "updated_at", "deleted_at"}).
			AddRow(uuid.New(), userID, "CAF", time.Now(), time.Now(), nil).
			AddRow(uuid.New(), userID, "AFG", time.Now(), time.Now(), nil)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE user_id = $1 AND deleted_at IS NULL`)).
			WithArgs(userID).
			WillReturnRows(rows)

		subs, err := repo.GetSubscriptionsByUserID(userID)

		assert.NoError(t, err)
		assert.Len(t, subs, 2)
		assert.Equal(t, "CAF", subs[0].CountryCode)
		assert.Equal(t, "AFG", subs[1].CountryCode)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - No subscriptions found", func(t *testing.T) {
		userID := uuid.New()

		rows := sqlmock.NewRows([]string{"id", "user_id", "country_code"})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE user_id = $1`)).
			WithArgs(userID).
			WillReturnRows(rows)

		subs, err := repo.GetSubscriptionsByUserID(userID)

		assert.NoError(t, err)
		assert.Empty(t, subs)
	})
}

func TestDeleteSubscription(t *testing.T) {
	db, mock, _ := SetupMockDB()
	repo := NewSubRepo(db)

	t.Run("Success - Soft delete subscription", func(t *testing.T) {
		userID := uuid.New()
		subID := uuid.New()

		rows := sqlmock.NewRows([]string{"id", "user_id", "country_code", "created_at", "updated_at", "deleted_at"}).
			AddRow(subID, userID, "USA", time.Now(), time.Now(), nil)

		mock.ExpectQuery(`SELECT \* FROM "subscriptions" WHERE id = \$1 AND user_id = \$2 AND deleted_at IS NULL.*`).
			WithArgs(subID, userID, 1). // GORM adds the '1' for the LIMIT clause
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "subscriptions" SET`)).
			WithArgs(
				userID,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				sqlmock.AnyArg(), // deleted_at
				subID,            // WHERE clause id
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeleteSubscription(userID, subID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failure - Subscription not found", func(t *testing.T) {
		userID := uuid.New()
		subID := uuid.New()

		mock.ExpectQuery(`SELECT \* FROM "subscriptions"`).
			WillReturnError(gorm.ErrRecordNotFound)

		err := repo.DeleteSubscription(userID, subID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error fetching subscription id")
	})
}
