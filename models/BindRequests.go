package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type BindRequest struct {
	gorm.Model
	ID uint `json:"id" gorm:"primaryKey"`

	ValidUntil    time.Time `json:"validUntil" gorm:"column:valid_until"`
	ServiceUserID string    `json:"serviceUserId" gorm:"index"`

	User2faID    string `json:"userId" gorm:"index;type:text;column:user_2fa_id"`
	Service2faID string `json:"serviceId" gorm:"index;type:text;column:service_2fa_id"`

	User2fa    *User2fa    `gorm:"foreignKey:User2faID;references:ID;constraint:OnDelete:CASCADE;"`
	Service2fa *Service2fa `gorm:"foreignKey:Service2faID;references:ID;constraint:OnDelete:CASCADE;"`
}

func (bindRequest *BindRequest) Create(db *gorm.DB) error {
	return db.Create(&bindRequest).Error
}

func GetBindRequestsByUserID(db *gorm.DB, userID string) []BindRequest {
	var bindRequests []BindRequest

	if err := db.Preload("User2fa").Preload("Service2fa").Find(&bindRequests).Where("user_2fa_id = ?", userID).Error; err != nil {
		return []BindRequest{}
	}

	return bindRequests
}

func GetAllBindRequests(db *gorm.DB) []BindRequest {
	var bindRequests []BindRequest

	if err := db.Preload("User2fa").Preload("Service2fa").Find(&bindRequests).Error; err != nil {
		return []BindRequest{}
	}

	return bindRequests
}

func GetBindRequestByID(db *gorm.DB, bindRequestID uint) (*BindRequest, error) {
	var bindRequest BindRequest

	if err := db.Preload("User2fa").Preload("Service2fa").First(&bindRequest, "id = ?", bindRequestID).Error; err != nil {
		return nil, err
	}

	return &bindRequest, nil
}

func GetLatestBindRequestForUser(db *gorm.DB, userID string, serviceID string) (*BindRequest, error) {
	var bindRequest BindRequest

	if err := db.Where("user_2fa_id = ? AND service_2fa_id = ? AND valid_until > ?", userID, serviceID, time.Now()).First(&bindRequest).Error; err != nil {
		return nil, err
	}

	if err := db.Where("user_2fa_id = ? AND service_2fa_id = ?", userID, serviceID).Delete(&BindRequest{}).Error; err != nil {
		return nil, err
	}

	return &bindRequest, nil
}

func DeleteBindRequestsForService(db *gorm.DB, userID string, serviceID string) error {
	var bindRequests []BindRequest

	if err := db.Where("user_2fa_id = ? AND service_2fa_id = ?", userID, serviceID).Find(&bindRequests).Error; err != nil {
		return err
	}

	if len(bindRequests) == 0 {
		return errors.New("bind request not found")
	}

	if err := db.Where("user_2fa_id = ? AND service_2fa_id = ?", userID, serviceID).Delete(&BindRequest{}).Error; err != nil {
		return err
	}

	return nil
}
