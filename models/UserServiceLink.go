package models

import "gorm.io/gorm"

type UserServiceLink struct {
	gorm.Model

	UserID    uint   `gorm:"index"`
	ServiceID uint   `gorm:"index"`
	UniqueKey string `gorm:"uniqueIndex:idx_user_service"`
}

func (UserServiceLink) TableName() string {
	return "user_service_link"
}
