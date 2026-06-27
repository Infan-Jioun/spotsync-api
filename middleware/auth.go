package middleware

import (
	"net/http"
	"os"
	"spotsync-api/dto"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Header থেকে token নাও
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: "Authorization header missing",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: "Invalid authorization format",
			})
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: "Invalid token claims",
			})
		}

		c.Set("userID", uint(claims["id"].(float64)))
		c.Set("userRole", claims["role"].(string))

		return next(c)
	}
}

// Admin only check
func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Get("userRole").(string)
		if role != "admin" {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Message: "Access denied. Admins only.",
			})
		}
		return next(c)
	}
}
