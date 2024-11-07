package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"selfhosted_2fa_sso/config"
	"selfhosted_2fa_sso/internal/db"
	"selfhosted_2fa_sso/internal/server"
)

func main() {
	config, err := config.LoadConfig()

	if err != nil {
		fmt.Println("Failed to load the config", err)
		os.Exit(1)
	}

	database, err := db.ConnectDatabase(config.Database.URL)

	if err != nil {
		fmt.Println("Failed to initialize the database:", err)
		os.Exit(1)
	}
	defer func() {
		sqlDB, _ := database.DB()
		sqlDB.Close()
	}()

	srv := server.NewServer(database, config)

	if err := srv.Start(); err != nil {
		fmt.Println("Failed to start the server:", err)
		os.Exit(1)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown
	fmt.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Error shutting down server:", err)
	} else {
		fmt.Println("Server stopped successfully")
	}
}
