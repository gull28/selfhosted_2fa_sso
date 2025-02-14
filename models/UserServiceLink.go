package models

import (
	"time"

	"gorm.io/gorm"
)

type UserServiceLink struct {
	gorm.Model

	ServiceUserID string `gorm:"index"` // unknown format
	ValidUntil    time.Time
	Enabled       bool

	User2faID    string `gorm:"index;type:text"`
	Service2faID string `gorm:"index"`

	User2fa    User2fa    `gorm:"foreignKey:User2faID;references:ID;constraint:OnDelete:CASCADE;"`
	Service2fa Service2fa `gorm:"foreignKey:Service2faID;constraint:OnDelete:CASCADE;"`
}

func (UserServiceLink) TableName() string {
	return "user_service_link"
}

func (usl *UserServiceLink) CreateUserServiceLinks(db *gorm.DB) error {
	return db.Create(usl).Error
}

func FetchUserServiceLinks(db *gorm.DB, userID string) ([]UserServiceLink, error) {
	var userServiceLinks []UserServiceLink
	if err := db.Preload("User2fa").Preload("Service2fa").Where("user_2fa_id = ?", userID).Find(&userServiceLinks).Error; err != nil {
		return nil, err
	}

	return userServiceLinks, nil
}

func (usl *UserServiceLink) IsUserAlreadyBound(db *gorm.DB) bool {
	var userServiceLink UserServiceLink

	result := db.Where("user2fa_id = ? AND service2fa_id = ?", usl.User2faID, usl.Service2faID).Find(&userServiceLink)
	if result.Error != nil {
		return false
	}

	return result.RowsAffected > 0
}

func IsAuthValid(db *gorm.DB, user2faID string, service2faID uint) (bool, error) {
	var userServiceLink UserServiceLink

	err := db.Where("service_user_id = ? AND service2fa_id = ?", user2faID, service2faID).Find(&userServiceLink).Error

	if err != nil {
		return false, err
	}

	return userServiceLink.ValidUntil.After(time.Now()), nil
}
