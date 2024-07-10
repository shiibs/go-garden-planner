package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shiibs/go-garden-planner/database"
	"github.com/shiibs/go-garden-planner/model"
)

type RelationshipData struct {
    PlantName   string `json:"plant_name"`
    FriendNames []string `json:"friend_names"`
    EnemyNames  []string `json:"enemy_names"`
}

func CreateRelationShipHandler(c *fiber.Ctx) error {
    var data RelationshipData
    if err := c.BodyParser(&data); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
    }

    // Find the plant ID corresponding to the provided plant name
    var plant model.Plant
    if err := database.DBConn.Where("name = ?", data.PlantName).First(&plant).Error; err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Plant not found"})
    }

    // Create relationships with friends
    for _, friendName := range data.FriendNames {
        var friend model.Plant
        if err := database.DBConn.Where("name = ?", friendName).First(&friend).Error; err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Friend plant not found"})
        }

        // Create friendship from plant to friend
        friendship := &model.Friend{
            PlantID:  plant.ID,
            FriendID: friend.ID,
        }
        if err := database.DBConn.Create(friendship).Error; err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to create friend relationship"})
        }

        // Create friendship from friend to plant (bidirectional)
        reverseFriendship := &model.Friend{
            PlantID:  friend.ID,
            FriendID: plant.ID,
        }
        if err := database.DBConn.Create(reverseFriendship).Error; err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to create reverse friend relationship"})
        }
    }

    // Create relationships with enemies
    for _, enemyName := range data.EnemyNames {
        var enemy model.Plant
        if err := database.DBConn.Where("name = ?", enemyName).First(&enemy).Error; err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Enemy plant not found"})
        }

        // Create enemyship from plant to enemy
        enemyship := &model.Enemy{
            PlantID:  plant.ID,
            EnemyID: enemy.ID,
        }
        if err := database.DBConn.Create(enemyship).Error; err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to create enemy relationship"})
        }

        // Create enemyship from enemy to plant (bidirectional)
        reverseEnemyship := &model.Enemy{
            PlantID:  enemy.ID,
            EnemyID: plant.ID,
        }
        if err := database.DBConn.Create(reverseEnemyship).Error; err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to create reverse enemy relationship"})
        }
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Relationships created successfully"})
}