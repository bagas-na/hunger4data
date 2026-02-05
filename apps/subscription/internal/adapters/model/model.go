package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	Id          uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserId      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;uniqueIndex:idx_user_country"`
	CountryCode string     `json:"country_code" gorm:"not null;uniqueIndex:idx_user_country"`
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
