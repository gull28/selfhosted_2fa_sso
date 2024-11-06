// cmd/api/main.go
package main

import (
	"log"
	"yourapp/config"
	"yourapp/internal/server"
	"yourapp/models"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	db, err := models.ConnectDatabase(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	s := server.NewServer(db, cfg)
	if err := s.Start(cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
