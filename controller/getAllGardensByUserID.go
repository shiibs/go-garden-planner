package controller

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/shiibs/go-garden-planner/database"
	"github.com/shiibs/go-garden-planner/model"
)

func GetUserData(c *fiber.Ctx) error{
	context := fiber.Map{
		"status": "OK",
		"msg": "All gardens route",
	}

	
	email := c.Locals("email")
	

	if email == nil {
		log.Println("Email not found")
		context["msg"] = "Email not found."
		return c.Status(fiber.StatusUnauthorized).JSON(context)
	}

	var user model.User

	if err := database.DBConn.Where("email = ?", email).First(&user).Error; err != nil {
		log.Println("User not found.")
		context["msg"] = "User not found."
		return c.Status(fiber.StatusBadRequest).JSON(context)
	}

	gardenLayout, err := GetAllGarden(user.ID)

	if err != nil {
        log.Println("failed to get gardens:", err)
    }


	// get garden details of the user if available
    gardens := make([]model.GardenDetails, len(gardenLayout))

    for _, garden := range gardenLayout {
        var data model.GardenDetails
        data.ID = garden.ID
        data.Name = garden.Name

        gardens = append(gardens, data)
    }

   context["userName"] = user.UserName
   context["gardens"] = gardens
   
   c.Status(200)
   return c.JSON(context)
} 

func GetAllGarden(userID uint) ([]model.GardenLayout, error){
    var gardens []model.GardenLayout

    result := database.DBConn.Where("user_id = ?", userID).Find(&gardens)

    return gardens, result.Error
}