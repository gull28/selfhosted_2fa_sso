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

func GetBindRequestByUserID(db *gorm.DB) (*BindRequest, error) {
	return nil, nil
}

func AcceptBindRequest(db *gorm.DB, userID string) error {
	return nil
}

func DeclineBindRequest(db *gorm.DB, userID string) error {
	return nil
}
