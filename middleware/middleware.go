package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/shiibs/go-garden-planner/auth"
)

func Authenticate(c *fiber.Ctx) error {
	cookie := c.Cookies("cookie")

	
	
	if cookie == "" {
		log.Println("cookie not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token not present."})
	}

	
	claims, msg := auth.ValidateToken(cookie)

	if msg != "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": msg})
		 
	}

	c.Locals("email", claims.Email)

	return c.Next()
}