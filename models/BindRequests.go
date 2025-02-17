package models

import (
	"time"

	"gorm.io/gorm"
)

type BindRequest struct {
	gorm.Model

	ValidUntil    time.Time
	ServiceUserID string `gorm:"index"`

	User2faID    string `gorm:"index;type:text;column:user_2fa_id"`
	Service2faID string `gorm:"index;type:text;column:service_2fa_id"`

	User2fa    *User2fa    `gorm:"foreignKey:User2faID;references:ID;constraint:OnDelete:CASCADE;"`
	Service2fa *Service2fa `gorm:"foreignKey:Service2faID;references:ID;constraint:OnDelete:CASCADE;"`
}

func (bindRequest *BindRequest) Create(db *gorm.DB) error {
	return db.Create(&bindRequest).Error
}

func GetBindRequestsByUserID(db *gorm.DB, userID string) ([]BindRequest, error) {
	var bindRequest []BindRequest

	if err := db.Preload("User2fa").Preload("Service2fa").Find(&bindRequest).Where("user_2fa_id = ? AND validUntil > ?", userID, time.Now()).Error; err != nil {
		return nil, err
	}

	return bindRequest, nil
}

func AcceptBindRequest(db *gorm.DB, userID string, serviceID string) error {
	// find /latest/ bind request and accept it if it exists and is still valid

	// delete it

	// create service link
}

func DeclineBindRequest(db *gorm.DB, userID string) error {
	// check if bind request entry exists

	// delete it
}
