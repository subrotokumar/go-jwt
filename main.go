package main

import (
	"go-jwt/controllers"
	"go-jwt/initializer"
	middlewares "go-jwt/middlewars"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	initializer.LoadENV()
	initializer.ConnectToDb()
}

func main() {
	app := fiber.New()
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})

	app.Post("/signup", controllers.Signup)
	app.Post("/login", controllers.LogIn)
	app.Post("/valid", middlewares.RequireAuth, controllers.ValidateToken)
	log.Fatal(app.Listen(":3000"))
}
