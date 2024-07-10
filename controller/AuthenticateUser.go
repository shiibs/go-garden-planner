package controller

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/shiibs/go-garden-planner/database"
)

func IsAuthenticated(c *fiber.Ctx) error {
    // Retrieve session
    sess, err := database.Store.Get(c)
    if err != nil {
        log.Println("Error retrieving session:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "statusText": "Error",
            "message":    "Session retrieval failed",
        })
    }

   
    defer sess.Save()

    // Check if user is authenticated
    userID := sess.Get("userID")
    userEmail := sess.Get("userEmail")
    log.Println("userID, email", userID, userEmail)
    if userID == nil || userEmail == nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "statusText": "Error",
            "message":    "Unauthorized",
        })
    }

    // Pass user ID and email to the next handler
    c.Locals("userID", userID)
    c.Locals("userEmail", userEmail)
    return c.Next()
}