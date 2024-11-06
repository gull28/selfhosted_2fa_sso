package server

import (
	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
}
