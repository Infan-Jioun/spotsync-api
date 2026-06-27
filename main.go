package main

import (
	"log"
	"os"
	"spotsync-api/config"
	"spotsync-api/handler"
	"spotsync-api/repository"
	"spotsync-api/service"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file load failed")
	}

	db := config.ConnectDB()

	// Dependency Injection — wire করো
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	api := e.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// Health check
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "SpotSync API is running 🚗",
		})
	})

	port := os.Getenv("PORT")
	log.Println("🚀 Server running on port", port)
	e.Logger.Fatal(e.Start(":" + port))
}
