package models

import "gorm.io/gorm"

type Service2fa struct {
	gorm.Model

	ID          uint `gorm:"uniqueIndex"`
	Name        string
	Description string
}

func (service *Service2fa) CreateService(db *gorm.DB) error {
	return db.Create(&service).Error
}

func GetServiceByID(db *gorm.DB, id string) (*Service2fa, error) {
	var service Service2fa
	err := db.First(&service, id).Error
	return &service, err
}
