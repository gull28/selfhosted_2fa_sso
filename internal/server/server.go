package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"selfhosted_2fa_sso/config"
	"selfhosted_2fa_sso/internal/ratelimit"
	"selfhosted_2fa_sso/middleware"
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

	ratelimiter := ratelimit.NewRateLimiter(10, 10)
	rateMiddleware := middleware.RateLimiterMiddleware(ratelimiter)

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	router.Use(rateMiddleware)

	router.Use(NoCacheStatic())

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

	routes.RegisterUserRoutes(s.router, s.db, s.config)

	sessionRoutes := s.router.Group("/session")
	sessionRoutes.Use(middleware.AuthMiddleware(s.db, s.config.JWT.Secret))
	routes.RegisterSessionRoutes(sessionRoutes, s.db, s.config)

	serviceRoutes := s.router.Group("/service")
	// serviceRoutes.Use(middleware.AuthMiddleware(s.db, s.config.JWT.Secret))
	routes.RegisterServiceRoutes(serviceRoutes, s.db)

	bindRoutes := s.router.Group("/bind")
	routes.RegisterBindRoutes(bindRoutes, s.db, s.config)
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

// dev hot reload css :PP
func NoCacheStatic() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, "/static/") {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}
		c.Next()
	}
}
