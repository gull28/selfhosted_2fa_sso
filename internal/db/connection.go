package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := Migrate(db); err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}

	fmt.Println("No error when connecting")
	return db, nil
}
