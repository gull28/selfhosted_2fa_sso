package models

import (
	"time"

	"gorm.io/gorm"
)

type UserServiceLink struct {
	gorm.Model

	ServiceUserID string `gorm:"index"` // unknown format
	ValidUntil    time.Time

	User2faID    uint `gorm:"index"`
	Service2faID uint `gorm:"index"`

	User2fa    User2fa    `gorm:"foreignKey:User2faID;constraint:OnDelete:CASCADE;"`
	Service2fa Service2fa `gorm:"foreignKey:Service2faID;constraint:OnDelete:CASCADE;"`
}

func (UserServiceLink) TableName() string {
	return "user_service_link"
}

func (usl *UserServiceLink) CreateUserServiceLink(db *gorm.DB) error {
	return db.Create(usl).Error
}

func (usc *UserServiceLink) IsUserAlreadyBound(db *gorm.DB) bool {
	var userServiceLink UserServiceLink

	result := db.Where("user2fa_id = ? AND service2fa_id = ?", usc.User2faID, usc.Service2faID).Find(&userServiceLink)
	if result.Error != nil {
		return false
	}

	return result.RowsAffected > 0
}

func IsAuthValid(db *gorm.DB, user2faID uint, service2faID uint) (bool, error) {
	var userServiceLink UserServiceLink

	result := db.Where("user2fa_id = ? AND service2fa_id = ?", user2faID, service2faID).Find(&userServiceLink)

	if result.Error != nil {
		return false, result.Error
	}

	if userServiceLink.ValidUntil.Before(time.Now()) {
		return false, nil
	}

	return true, nil
}
