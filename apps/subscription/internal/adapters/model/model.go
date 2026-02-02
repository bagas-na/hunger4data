package model

type Subscription struct {
	Id         int64 `gorm:"primaryKey"`
	Id_user    int64 `json:"id_user" gorm:"not null"`
	Id_country int64 `json:"id_country" gorm:"not null"`
}

type Country struct {
	Id   int64  `gorm:"primaryKey"`
	Name string `json:"location_name" gorm:"not null"`
}
