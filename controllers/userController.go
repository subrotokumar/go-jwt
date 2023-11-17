package controllers

import (
	"go-jwt/initializer"
	"go-jwt/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *fiber.Ctx) error {
	var body struct {
		Email    string
		Password string
	}

	if c.BodyParser(&body) != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}
	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializer.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User successfully created",
	})
}

func LogIn(c *fiber.Ctx) error {
	var body struct {
		Email    string
		Password string
	}

	if c.BodyParser(&body) != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	var user models.User
	initializer.DB.Find(&user, "email = ?", body.Email)
	if user.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Invalid email/password",
		})
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Invalid email/password",
		})
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, error := token.SignedString([]byte(os.Getenv("SECRET")))
	if error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Faild To create token",
		})
	}

	cookie := fiber.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure:   true,
	}
	c.Cookie(&cookie)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": tokenString,
	})
}

func ValidateToken(c *fiber.Ctx) error {
	user := c.Get("user", "error")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": user,
	})
}
