package controller

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/shiibs/go-garden-planner/database"
	"github.com/shiibs/go-garden-planner/model"
)

func DeleteGarden(c *fiber.Ctx) error {
	context := fiber.Map{
		"status": "OK",
		"msg": "Garden Deleted succesfully",
	}

	id := c.Params("id")

	record := new(model.GardenLayout)

	database.DBConn.First(&record, id)

	if record.ID == 0 {
		log.Println("Record not found")
		context["msg"] = "Record Not Found"
		c.Status(400)
		return c.JSON(context)
	}
	
	result := database.DBConn.Delete(record)

	if result.Error != nil {
		log.Println("Something went wrong")
		context["msg"] = "something went wrong"
		c.Status(400)
		return c.JSON(context)
	}

	c.Status(200)
	return c.JSON(context)
}