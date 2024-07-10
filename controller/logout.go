package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
)


func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "cookie",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})
	return c.SendStatus(fiber.StatusOK)
}