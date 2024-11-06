package models

import "gorm.io/gorm"

type Service2FA struct {
	gorm.Model

	ID          string `gorm:"uniqueIndex"`
	Name        string
	Description string
}

func (service *Service2FA) CreateService(db *gorm.DB) error {
	return db.Create(u).Error
}

func GetServiceByID(db *gorm.DB, id uint) (*Service2FA, error) {
	var service Service2FA
	err := db.First(&service, id).Error
	return &service, err
}
