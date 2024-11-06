package main

import (
	"log"
	"selfhosted_2fa_sso/config"
)

func main() {
	appConfig, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// services

	log.Printf("Running on port: %s", appConfig.App.Port)
	log.Printf("Connecting to database: %s", appConfig.Database.URL)
}
