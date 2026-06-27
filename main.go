package main

import (
	"log"
	"os"
	"spotsync-api/config"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file load failed")
	}

	db := config.ConnectDB()
	_ = db

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "SpotSync API is running ",
		})
	})

	// Server start
	port := os.Getenv("PORT")
	log.Println(" Server running on port", port)
	e.Logger.Fatal(e.Start(":" + port))
}
