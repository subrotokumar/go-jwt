package middlewares

import (
	"fmt"
	"go-jwt/initializer"
	"go-jwt/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *fiber.Ctx) error {
	tokenString := c.Cookies("authorization")
	if tokenString == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		var user models.User
		initializer.DB.Find(&user, claims["email"])
		if user.Id == 0 {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		c.Set("user", user.Email)
		c.Next()

		return c.Next()
	} else {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
}
