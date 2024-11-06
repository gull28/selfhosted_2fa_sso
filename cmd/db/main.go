package db

import (
	"log"
	"selfhosted_2fa_sso/config"
)

func main() {
	appConfig, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("AppConfig error %v", err)
	}
}
