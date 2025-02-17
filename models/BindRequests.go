package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type BindRequest struct {
	gorm.Model

	ValidUntil    time.Time `gorm:"column:valid_until"`
	ServiceUserID string    `gorm:"index"`

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

	if err := db.Preload("User2fa").Preload("Service2fa").Find(&bindRequest).Where("user_2fa_id = ? AND valid_until > ?", userID, time.Now()).Error; err != nil {
		return nil, err
	}

	return bindRequest, nil
}

func DeleteBindRequestsForService(db *gorm.DB, userID string, serviceID string) error {
	var bindRequests []BindRequest

	if err := db.Where("user_2fa_id = ? AND service_2fa_id = ? AND valid_until > ?", userID, serviceID, time.Now()).Find(&bindRequests).Error; err != nil {
		return err
	}

	if len(bindRequests) == 0 {
		return errors.New("bind request not found")
	}

	if err := db.Where("user_2fa_id = ? AND service_2fa_id = ?", userID, serviceID).Delete(&BindRequest{}).Error; err != nil {
		return err
	}

	return nil
}
