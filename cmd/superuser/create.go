package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"selfhosted_2fa_sso/config"
	"selfhosted_2fa_sso/internal/db"
	"selfhosted_2fa_sso/models"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	cfg, err := config.LoadConfig()

	if err != nil {
		fmt.Printf("Error loading the config: %v\n", err)
		os.Exit(1)
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

		if len(username) > 20 {
			fmt.Println("Username too long, keep it under 20 characters!")
			continue
		}

		if len(username) < 3 {
			fmt.Println("Username too short, keep it above 2 characters!")
			continue
		}

		var existingUser models.SuperUser
		if err := database.Where("username = ?", username).First(&existingUser).Error; err == nil {
			fmt.Println("Username already taken. Please enter a different username.")
		} else {
			break
		}
	}

	var password string
	for {
		fmt.Print("Enter superuser password: ")
		password, _ = reader.ReadString('\n')
		password = strings.TrimSpace(password)

		if len(password) < 6 {
			fmt.Println("Password must be at least 6 characters long")
			continue
		}

		if len(password) >= 128 {
			fmt.Println("Password must be less than 128 characters long")
			continue
		}

		break
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Failed to hash password: %v\n", err)
		os.Exit(1)
	}

	superuser := models.SuperUser{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	validate := validator.New()
	if err := validate.Struct(&superuser); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		fmt.Println("Try again")
		return
	}

	if err := database.Create(&superuser).Error; err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Superuser created successfully!")
}
