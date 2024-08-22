package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shiibs/go-garden-planner/auth"
	"github.com/shiibs/go-garden-planner/controller"
	"github.com/shiibs/go-garden-planner/middleware"
)

func SetupRouters(app *fiber.App) {
	app.Get("/get_plants", controller.GetPlantList)
	app.Post("/add_plant", controller.CreatePlant)
	app.Post("/add_relation", controller.CreateRelationShipHandler)
	app.Get("/auth/login", auth.GoogleLoginHandler)
	app.Get("/auth/callback", auth.GoogleCallbackHandler)
	// app.Get("/getUser", auth.GetUserHandler)

	api := app.Group("/api")

	api.Use(middleware.Authenticate)

	api.Post("/create_gardenLayout", controller.PostGardenPlanner)
	api.Get("/garden_layout/:id", controller.GetGardenLayoutWithID)
	api.Get("/refreshToken", controller.RefreshToken)
	api.Delete("/delete_garden/:id", controller.DeleteGarden)
	api.Get("/get_user_data", controller.GetUserData)

	api.Post("/logout", controller.Logout)

	app.Static("/", "./dist")
}
