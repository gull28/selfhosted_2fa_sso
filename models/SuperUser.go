package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// This has NO connection to User2fa, this is not a 2fa user
type SuperUser struct {
	gorm.Model
	Username     string `gorm:"unique;not null;size:20" validate:"required,min=3,max=20"`
	PasswordHash string `gorm:"not null"`
}

func FindSuperUserByUsername(db *gorm.DB, username string) (*SuperUser, error) {
	var user SuperUser
	err := db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func Create(db *gorm.DB, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := SuperUser{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}
	return db.Create(&user).Error
}

func FindSuperUserByID(db *gorm.DB, id uint) (*SuperUser, error) {
	var user SuperUser
	err := db.First(&user, id).Error

	return &user, err
}
