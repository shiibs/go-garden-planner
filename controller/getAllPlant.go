package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shiibs/go-garden-planner/database"
	"github.com/shiibs/go-garden-planner/model"
)

func GetPlantList(c *fiber.Ctx) error {
	context := fiber.Map {
		"statusText": "OK",
		"message": "Plant list",
	}
	var plantList []model.Plant

	if err := database.DBConn.Find(&plantList).Error; err != nil {
		context["statusText"] = "Error"
		context["message"] = "Failed to fetch data from database"
		return c.Status(fiber.StatusInternalServerError).JSON(context)
	} 

	 for i := range plantList {
		var enemies []model.Enemy
		var friends []model.Friend
		database.DBConn.Where("plant_id = ?", plantList[i].ID).Find(&enemies)
		database.DBConn.Where("plant_id = ?", plantList[i].ID).Find(&friends)
		plantList[i].EnemyPlants = enemies
		plantList[i].FriendPlants = friends
	 }

	context["plant_list"] = plantList

	
	
	return c.Status(200).JSON(context)
}