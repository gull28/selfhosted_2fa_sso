package models

import "gorm.io/gorm"

type Service2FA struct {
	gorm.Model

	ServiceID   string `gorm:"uniqueIndex"`
	Name        string
	Description string
}
