package models

import "gorm.io/gorm"

type UserServiceLink struct {
	gorm.Model

	UserID    uint `gorm:"index"` // outer service user id
	ServiceID uint `gorm:"index"`

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
