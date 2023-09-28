package server

import (
	"authx/controllers"
	"authx/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupAndListen() {

	router := fiber.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	router.Get("/api", index)

	router.Post("/api/register", controllers.Register)
	router.Post("/api/login", controllers.Login)
	router.Post("/api/logout", controllers.Logout)
	router.Get("/api/users", middleware.VerifyToken, controllers.FetchAllUsers)
	router.Get("/api/search/:term", middleware.VerifyToken, controllers.Search)
	router.Post("/api/photo", middleware.VerifyToken, controllers.UploadPhoto)
	router.Patch("/api/update", middleware.VerifyToken, controllers.UpdateUser)

	router.Listen(":3000")

}

func index(
	c *fiber.Ctx) error {
	return c.SendString("Welcome to Auth App")
}
