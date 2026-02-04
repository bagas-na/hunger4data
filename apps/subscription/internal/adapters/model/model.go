package model

import "github.com/google/uuid"

type Subscription struct {
	Id         uuid.UUID `gorm:"type:uuid;primaryKey"`
	Id_user    uuid.UUID `json:"id_user" gorm:"type:uuid;not null"`
	Id_country uuid.UUID `json:"id_country" gorm:"type:uuid;not null"`
}

type Country struct {
	Id   int64  `gorm:"primaryKey"`
	Name string `json:"location_name" gorm:"not null"`
}
