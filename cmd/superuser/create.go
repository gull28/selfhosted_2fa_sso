package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"selfhosted_2fa_sso/config"
	"selfhosted_2fa_sso/internal/db"
	"selfhosted_2fa_sso/models"
)

//nolint:unused
func main() {
	reader := bufio.NewReader(os.Stdin)
	cfg, err := config.LoadConfig()

	if err != nil {
		fmt.Printf("Error loading the config %v", err)
	}

	database, err := db.ConnectDatabase(cfg.Database.URL)

	if err != nil {
		fmt.Println("Failed to initialize the database:", err)
		os.Exit(1)
	}
	defer func() {
		sqlDB, _ := database.DB()
		sqlDB.Close()
	}()

	var username string
	for {
		fmt.Print("Enter superuser username: ")
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username)

		var existingUser models.SuperUser
		if err := database.Where("username = ?", username).First(&existingUser).Error; err == nil {
			fmt.Println("Username already taken. Please enter a different username.")
		} else {
			break
		}
	}

	fmt.Print("Enter superuser password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Printf("Failed to hash password: %v\n", err)
		os.Exit(1)
	}

	superuser := models.SuperUser{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	err = superuser.Create(database)

	if err != nil {
		fmt.Printf("Error creating user %v", err)
		os.Exit(1)
	}
}
