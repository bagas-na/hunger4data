package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Run migrations
	db.AutoMigrate(&User{}, &Payment{}, &PaymentEvent{})
	return db
}

func TestCreatePayment_Transaction(t *testing.T) {
	gormDB := setupTestDB(t)
	repo := NewPaymentRepo(gormDB)
	ctx := context.Background()

	userID := uuid.New()
	p := &Payment{
		UserID:      userID,
		Amount:      1000,
		Currency:    "USD",
		Status:      StatusPending,
		CountryCode: "USA",
	}

	created, err := repo.CreatePayment(ctx, p)

	// Assert
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, created.ID)

	// Verify PaymentEvent was also created
	var event PaymentEvent
	err = gormDB.First(&event, "payment_id = ?", created.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, EventCreated, event.EventType)
}

func TestFindPaymentByID_WithPreload(t *testing.T) {
	gormDB := setupTestDB(t)
	repo := NewPaymentRepo(gormDB)

	// Seed a user and a payment
	user := User{Id: uuid.New(), Username: "tester"}
	gormDB.Create(&user)

	pay := Payment{ID: uuid.New(), UserID: user.Id, Status: StatusPaid}
	gormDB.Create(&pay)

	// Act
	found, err := repo.FindPaymentByID(context.Background(), pay.ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "tester", found.User.Username)
}
