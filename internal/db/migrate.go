package db

import (
	"selfhosted_2fa_sso/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Service2fa{}, &models.User2fa{}, &models.UserServiceLink{})
}
