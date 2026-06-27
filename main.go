package main

import (
	"log"
	"os"
	"spotsync-api/config"
	"spotsync-api/handler"
	appMiddleware "spotsync-api/middleware"
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

	// Repositories
	userRepo := repository.NewUserRepository(db)
	zoneRepo := repository.NewZoneRepository(db)

	// Services
	authService := service.NewAuthService(userRepo)
	zoneService := service.NewZoneService(zoneRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	zoneHandler := handler.NewZoneHandler(zoneService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	api := e.Group("/api/v1")

	// ✅ Public auth routes
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// ✅ Public zone routes (GET only)
	api.GET("/zones", zoneHandler.GetAll)
	api.GET("/zones/:id", zoneHandler.GetByID)

	// ✅ Protected zone routes (JWT + Admin)
	api.POST("/zones", zoneHandler.Create, appMiddleware.JWTMiddleware, appMiddleware.AdminOnly)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "SpotSync API is running 🚗",
		})
	})

	port := os.Getenv("PORT")
	log.Println("🚀 Server running on port", port)
	e.Logger.Fatal(e.Start(":" + port))
}
