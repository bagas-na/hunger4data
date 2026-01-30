package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Payment{}); err != nil {
		fmt.Println("payment automigration failed")
	}
	fmt.Println("payment automigration complete")

	if err := db.AutoMigrate(&PaymentEvent{}); err != nil {
		fmt.Println("payment_event automigration failed")
	}
	fmt.Println("payment_event automigration complete")

	return db, nil
}
