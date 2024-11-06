package models

import (
	"gorm.io/gorm", 

	"github.com/pquerna/otp/totp"
)

type User2FA struct {
	gorm.Model
	UUID       string `gorm:"uniqueIndex" json:"uuid"`
	Username   string `gorm:"unique" json:"username" binding:"required,min=3,max=20"`
	TOTPSecret string `gorm:"not null" json:"totp_secret" binding:"required"`
}

func (u *User2FA) CreateUser(db *gorm.DB) error {
	return db.Create(u).Error
}

func GetUserByID(db *gorm.DB, id uint) (*User2FA, error) {
	var user User2FA
	err := db.First(&user, id).Error
	return &user, err
}

func (u *User2FA) UpdateUser(db *gorm.DB) error {
	return db.Save(u).Error
}

func DeleteUser(db *gorm.DB, id uint) error {
	return db.Delete(&User2FA{}, id).Error
}
