package service

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	Id         int            `gorm:"primaryKey" json:"id"`
	Username   string         `gorm:"unique;not null" json:"username"`
	Password   string         `gorm:"not null" json:"password"`
	Role       string         `gorm:"type:varchar(20);default:'user'" json:"role"`
	Created_At time.Time      `json:"created_at"`
	Updated_At time.Time      `json:"updated_at"`
	Deleted_At gorm.DeletedAt `gorm:"index" json:"-"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
