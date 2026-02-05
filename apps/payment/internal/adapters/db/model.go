package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Username         string         `gorm:"unique;not null" json:"username"`
	PasswordHash     string         `gorm:"not null" json:"password_hash"`
	Role             string         `gorm:"type:varchar(20);default:'user'" json:"role"`
	IsActivated      bool           `gorm:"not null;default:false" json:"is_activated"`
	ActivationString string         `gorm:"not null" json:"activation_string"`
	Created_At       time.Time      `json:"created_at"`
	Updated_At       time.Time      `json:"updated_at"`
	Deleted_At       gorm.DeletedAt `gorm:"index" json:"-"`
}

type Payment struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	User   User      `gorm:"foreignKey:UserID"`

	CountryCode string `gorm:"not null"`

	TransactionType string `gorm:"size:32;not null"` // payment, refund

	Amount   int64  `gorm:"not null"`
	Currency string `gorm:"size:3;not null"`

	Provider          string `gorm:"size:32;not null"`
	ProviderSessionID string `gorm:"size:255"`

	Status PaymentStatus `gorm:"varchar(32);not null"` // pending, paid, failed, expired

	CreatedAt time.Time
	UpdatedAt time.Time
}

type PaymentEvent struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	PaymentID uuid.UUID `gorm:"type:uuid;not null"`
	Payment   Payment   `gorm:"foreignKey:PaymentID"`

	EventType PaymentEventType `gorm:"size:32;not null"` // created, pending, paid, failed, expired
	Source    string           `gorm:"size:16;not null"` // app, webhook

	ProviderEventID string //`gorm:"uniqueIndex"`

	CreatedAt time.Time
}

type PaymentStatus string

const (
	StatusCreated PaymentStatus = "created"
	StatusPending PaymentStatus = "pending"
	StatusPaid    PaymentStatus = "paid"
	StatusFailed  PaymentStatus = "failed"
	StatusExpired PaymentStatus = "expired"
)

type PaymentEventType string

const (
	EventCreated PaymentEventType = "created"
	EventPending PaymentEventType = "pending"
	EventPaid    PaymentEventType = "paid"
	EventFailed  PaymentEventType = "failed"
	EventExpired PaymentEventType = "expired"
)
