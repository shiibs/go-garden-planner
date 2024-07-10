package controller

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/shiibs/go-garden-planner/database"
	"github.com/shiibs/go-garden-planner/model"
)

func CreatePlant(c *fiber.Ctx) error {
	context := fiber.Map {
		"statusText": "OK",
		"message": "Plant Added",
	}

	plant := new(model.Plant)

	if err := c.BodyParser(plant); err != nil {
		log.Panicln("Error in parsing request")
		context["statusText"] = "Error"
		context["message"] = "Error in parsing request"
		return c.Status(fiber.StatusInternalServerError).JSON(context)
	}

	resutl := database.DBConn.Create(plant)

	if resutl.Error != nil {
		log.Println("Error in saving data")
		context["statusText"] = "Error"
		context["message"] = "Error in saving data"
		return c.Status(fiber.StatusInternalServerError).JSON(context)
	}

	context["data"] = plant
	c.Status(201)
	return c.JSON(context)
}