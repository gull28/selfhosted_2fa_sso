package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User2fa struct {
	ID         string `gorm:"primaryKey"`
	Username   string `gorm:"unique" json:"username" binding:"required,min=3,max=20"`
	TOTPSecret string `gorm:"not null" json:"totp_secret" binding:"required"`
}

func (u *User2fa) CreateUser(db *gorm.DB) error {
	return db.Create(&u).Error
}

func GetUserByID(db *gorm.DB, id uint) (*User2fa, error) {
	var user User2fa
	err := db.First(&user, id).Error
	return &user, err
}

func (u *User2fa) UpdateUser(db *gorm.DB) error {
	return db.Save(u).Error
}

func DeleteUser(db *gorm.DB, id uint) error {
	return db.Delete(&User2fa{}, id).Error
}

func (u *User2fa) BeforeCreate(_ *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}
