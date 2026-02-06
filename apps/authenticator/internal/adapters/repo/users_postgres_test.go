package repo

import (
	"errors"
	"regexp"
	"testing"

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

func TestCreateUser(t *testing.T) {
	db, mock, _ := SetupMockDB()
	repo := NewUserRepo(db)

	t.Run("success", func(t *testing.T) {
		user := User{
			Id:           uuid.New(),
			Username:     "test@gmail.com",
			PasswordHash: "securepassword",
			Role:         "user",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.CreateUser(user)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate username", func(t *testing.T) {
		user := User{
			Id:       uuid.New(),
			Username: "existing_user",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
			WillReturnError(&pgconn.PgError{
				Code: "23505",
			})
		mock.ExpectRollback()

		err := repo.CreateUser(user)

		assert.ErrorIs(t, err, gorm.ErrDuplicatedKey)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("generic db error", func(t *testing.T) {
		user := User{
			Id:       uuid.New(),
			Username: "fail_user",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
			WillReturnError(errors.New("db down"))
		mock.ExpectRollback()

		err := repo.CreateUser(user)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUser(t *testing.T) {
	db, mock, _ := SetupMockDB()
	repo := NewUserRepo(db)

	t.Run("success - activated user found", func(t *testing.T) {
		username := "johndoe"
		rows := sqlmock.NewRows([]string{"id", "username", "is_activated", "deleted_at"}).
			AddRow(uuid.New(), username, true, nil) // State: true

		// Note the trailing space in "deleted_at IS NULL " to match your repo code
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "users" WHERE (username = $1 AND is_activated = $2 AND deleted_at IS NULL ) AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $3`)).
			WithArgs(username, true, 1). // We expect the repo to pass 'true' here
			WillReturnRows(rows)

		user, err := repo.GetByUsername(username)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.True(t, user.IsActivated)
	})

	t.Run("failure - user exists but not activated", func(t *testing.T) {
		username := "unactivated_user"

		// The DB returns nothing because the SQL includes "is_activated = true"
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "users" WHERE (username = $1 AND is_activated = $2 AND deleted_at IS NULL )`)).
			WithArgs(username, true, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"})) // Empty result

		user, err := repo.GetByUsername(username)

		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		assert.Nil(t, user)
	})
}

func TestUpdateUser(t *testing.T) {
	db, mock, _ := SetupMockDB()
	repo := NewUserRepo(db)

	t.Run("success", func(t *testing.T) {
		username := "johndoe"
		updated := User{
			Role: "admin",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
			WithArgs(
				updated.Role,
				username,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.UpdateUser(username, updated)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteUser(t *testing.T) {
	db, mock, _ := SetupMockDB()
	repo := NewUserRepo(db)

	t.Run("success", func(t *testing.T) {
		username := "johndoe"

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "deleted_at"`)).
			WithArgs(
				sqlmock.AnyArg(),
				username,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.DeleteUser(username)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
