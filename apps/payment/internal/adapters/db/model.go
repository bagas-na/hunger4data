package db

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	CountryID uuid.UUID `gorm:"type:uuid;not null"`

	TransactionType string `gorm:"size:32;not null"` // payment, refund

	Amount   int64  `gorm:"not null"`
	Currency string `gorm:"size:3;not null"`

	Provider          string `gorm:"size:32;not null"`
	ProviderSessionID string `gorm:"size:255"`

	Status string `gorm:"varchar(32);not null"` // pending, paid, failed, expired

	CreatedAt time.Time
	UpdatedAt time.Time
}

type PaymentEvent struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	PaymentID uuid.UUID `gorm:"type:uuid;not null"`
	Payment   Payment   `gorm:"foreignKey:PaymentID"`

	EventType string `gorm:"size:32;not null"` // created, pending, paid, failed, expired
	Source    string `gorm:"size:16;not null"` // app, webhook

	ProviderEventID string `gorm:"uniqueIndex"`

	CreatedAt time.Time
}
