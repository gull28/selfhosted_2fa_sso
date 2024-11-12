package models

import "gorm.io/gorm"

// This has NO connection to User2fa, this is not a 2fa user
type SuperUser struct {
	gorm.Model

	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null;size:20" validate:"required,min=3,max=20"`
	PasswordHash string `gorm:"not null"`
}

func FindSuperUserRegistered(db *gorm.DB) (*SuperUser, error) {
	var user SuperUser
	err := db.First(&user).Error

	return &user, err
}

func FindSuperUserByID(db *gorm.DB, id uint) (*SuperUser, error) {
	var user SuperUser
	err := db.First(&user, id).Error

	return &user, err
}

func (su *SuperUser) Create(db *gorm.DB) error {
	return db.Create(&su).Error
}
