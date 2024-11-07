package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"selfhosted_2fa_sso/config"
	"selfhosted_2fa_sso/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	db         *gorm.DB
	router     *gin.Engine
	config     *config.Config
	httpServer *http.Server
}

func NewServer(db *gorm.DB, cfg *config.Config) *Server {
	router := gin.Default()

	s := &Server{
		db:     db,
		router: router,
		config: cfg,
	}

	s.setupRoutes()

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return s
}

func (s *Server) setupRoutes() {
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	routes.RegisterUserRoutes(s.router, s.db)
	routes.RegisterServiceRoutes(s.router, s.db)
}

func (s *Server) Start() error {
	log.Printf("Starting server on port %s\n", s.config.App.Port)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
