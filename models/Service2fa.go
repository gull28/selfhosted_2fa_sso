package models

import "gorm.io/gorm"

type Service2fa struct {
	gorm.Model

	ID          uint   `json:"serviceId" gorm:"uniqueIndex"`
	Name        string `json:"name" gorm:"unique;not null"`
	Description string `json:"description" gorm:"not null"`
}

func (service *Service2fa) Create(db *gorm.DB) error {
	return db.Create(&service).Error
}

func GetServiceByID(db *gorm.DB, id string) (*Service2fa, error) {
	var service Service2fa
	err := db.First(&service, id).Error
	return &service, err
}
