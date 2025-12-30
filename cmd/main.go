package main

import (
	"database/sql"
	"e-ticketing/config"
	"e-ticketing/internal/handler"
	"e-ticketing/internal/repository"
	"e-ticketing/internal/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("✅ Database connected successfully")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	emailSvc := service.NewEmailService(cfg)
	whatsappSvc := service.NewWhatsAppService(cfg)
	authSvc := service.NewAuthService(userRepo, emailSvc, whatsappSvc, cfg)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authSvc)

	// Setup Gin router
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "E-Ticketing API is running"})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/select-verification", authHandler.SelectVerificationMethod)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
			auth.GET("/verify-email", authHandler.VerifyEmailToken)
			auth.POST("/resend-otp", authHandler.ResendOTP)
		}
	}

	// Start server
	log.Printf("🚀 Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}