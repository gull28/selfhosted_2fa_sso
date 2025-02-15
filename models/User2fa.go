package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User2fa struct {
	ID         string `gorm:"primaryKey;type:text" json:"id"`
	Username   string `gorm:"unique" json:"username" binding:"required,min=3,max=20"`
	TOTPSecret string `gorm:"not null" json:"totp_secret" binding:"required"`
}

func (u *User2fa) CreateUser(db *gorm.DB) error {
	return db.Create(&u).Error
}

func GetUserByID(db *gorm.DB, id string) (*User2fa, error) {
	var user User2fa
	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(db *gorm.DB, username string) (*User2fa, error) {
	var user User2fa

	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func FetchAllUsers(db *gorm.DB) ([]*User2fa, error) {
	var users []*User2fa
	err := db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
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
