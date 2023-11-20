package routes

import (
	"github.com/alireza-frj4/BlogBackEnd/controller"
	"github.com/alireza-frj4/BlogBackEnd/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controller.Register)
	app.Post("/api/login", controller.Login)

	app.Use(middleware.IsAuthenticate)
	app.Post("/api/post", controller.CreatePost)
	app.Post("/api/upload-images", controller.Upload)
	app.Get("/api/allpost", controller.AllPost)
	app.Get("/api/allpost/:id", controller.DetailPost)
	app.Get("/api/uniquepost", controller.UniqePost)
	app.Delete("/api/deletepost/:id", controller.DeletePost)
	app.Static("api/uploads", "./uploads")
}
