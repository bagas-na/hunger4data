package db

import (
	"context"
	"fmt"

	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type NotificationLog struct {
	ID        uuid.UUID `gorm:"primaryKey;autoIncrement"`
	ToEmail   string    `gorm:"not null"`
	Subject   string    `gorm:"not null"`
	Status    string    `gorm:"not null"` // sent | failed
	Error     *string   `gorm:""`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func ConnectDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&NotificationLog{}); err != nil {
		fmt.Println("notification automigration failed")
	}
	fmt.Println("notification automigration complete")

	return db, nil
}

type NotificationLogRepo struct {
	db *gorm.DB
}

func NewNotificationLogRepo(db *gorm.DB) *NotificationLogRepo {
	return &NotificationLogRepo{db: db}
}

func (r *NotificationLogRepo) Create(ctx context.Context, log *NotificationLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
