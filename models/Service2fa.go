package models

import (
	"time"

	"gorm.io/gorm"
)

type Service2fa struct {
	ID          uint   `json:"serviceId" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"unique;not null;size=20" binding:"required,min=3,max=20"`
	Description string `json:"description" binding:"max=512" gorm:"size:512"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (service *Service2fa) Create(db *gorm.DB) error {
	return db.Create(&service).Error
}

func GetServiceByID(db *gorm.DB, id string) (*Service2fa, error) {
	var service Service2fa
	err := db.First(&service, id).Error

	return &service, err
}

func GetAllServices(db *gorm.DB) ([]*Service2fa, error) {
	var services []*Service2fa
	if err := db.Find(&services).Error; err != nil {
		return nil, err
	}

	return services, nil
}

func DeleteService(db *gorm.DB, id uint) error {
	return db.Delete(&Service2fa{}, id).Error
}
