package repo

import (
	"errors"
	"regexp"
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

func TestCreateUser(t *testing.T) {
	db, mock, _ := SetupMockDB()
	repo := NewUserRepo(db)
	t.Run("Successful creation", func(t *testing.T) {
		userID := uuid.New()
		user := User{
			Id:           userID,
			Username:     "test@gmail.com",
			PasswordHash: "securepassword",
			Role:         "user",
		}
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
			WithArgs(
				user.Id,
				user.Username,
				user.PasswordHash,
				user.Role,
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				sqlmock.AnyArg(), // DeletedAt
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.CreateUser(user)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Failed creation", func(t *testing.T) {
		userID := uuid.New()
		user := User{
			Id:           userID,
			Username:     "test@gmail.com",
			PasswordHash: "securepassword",
			Role:         "user",
		}
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
			WithArgs(sqlmock.AnyArg(), "test@gmail.com", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("duplicate key value violates unique constraint"))

		mock.ExpectRollback()

		err := repo.CreateUser(user)

		assert.Error(t, err)
	})
	t.Run("Duplicate username error", func(t *testing.T) {
		user := User{
			Id:       uuid.New(),
			Username: "existing_user",
		}

		pgErr := &pgconn.PgError{
			Code:    "23505",
			Message: "duplicate key value violates unique constraint",
		}

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
			WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(pgErr)

		mock.ExpectRollback()

		err := repo.CreateUser(user)

		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrDuplicatedKey)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUser(t *testing.T) {
	db, mock, _ := SetupMockDB()
	repo := NewUserRepo(db)
	t.Run("Successful Get", func(t *testing.T) {
		username := "johndoe"
		expectedUser := User{
			Id:       uuid.New(),
			Username: username,
			Role:     "user",
		}
		rows := sqlmock.NewRows([]string{"id", "username", "password", "role", "created_at", "updated_at", "deleted_at"}).
			AddRow(expectedUser.Id, expectedUser.Username, "hashed_pw", expectedUser.Role, time.Now(), time.Now(), nil)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(username, 1).
			WillReturnRows(rows)

		user, err := repo.GetByUsername(username)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, username, user.Username)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed Get", func(t *testing.T) {
		username := "nonexistent"

		// Return empty rows to simulate "not found"
		rows := sqlmock.NewRows([]string{"id", "username", "role"})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1`)).
			WithArgs(username, 1).
			WillReturnRows(rows)

		user, err := repo.GetByUsername(username)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, gorm.ErrRecordNotFound, err) // Verify it returns GORM's specific error
	})
}

func TestUpdateUser(t *testing.T) {
	db, mock, _ := SetupMockDB()
	repo := NewUserRepo(db)
	t.Run("Successful update", func(t *testing.T) {
		username := "johndoe"
		updatedData := User{
			Role:       "admin",
			Updated_At: time.Now(),
		}
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "role"=$1,"updated_at"=$2 WHERE username = $3 AND "users"."deleted_at" IS NULL`)).
			WithArgs(
				updatedData.Role,
				updatedData.Updated_At,
				username,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		err := repo.UpdateUser(username, updatedData)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestDeleteUser(t *testing.T) {
	db, mock, _ := SetupMockDB()
	repo := NewUserRepo(db)

	t.Run("Successful soft delete", func(t *testing.T) {
		username := "johndoe"

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "deleted_at"=`)).
			WithArgs(
				sqlmock.AnyArg(),
				username,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeleteUser(username)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}
