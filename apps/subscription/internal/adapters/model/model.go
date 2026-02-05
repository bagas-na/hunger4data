package model

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

type Subscription struct {
	Id     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	User   User      `gorm:"foreignKey:UserID"`

	CountryCode string     `json:"country_code" gorm:"not null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"not null"`
	UpdateAt    time.Time  `json:"updated_at" gorm:"not null"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type Country struct {
	Id                        int64   `gorm:"primaryKey"`
	Name                      string  `json:"location_name" gorm:"not null"`
	LocationCode              string  `json:"location_code" gorm:"not null"`
	IpcPhase                  string  `json:"ipc_phase,omitempty"`
	PopulationInPhase         int64   `json:"population_in_phase,omitempty"`
	PopulationFractionInPhase float64 `json:"population_fraction_in_phase,omitempty"`
}
