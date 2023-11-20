package main

import (
	"log"
	"os"

	"github.com/alireza-frj4/BlogBackEnd/database"
	"github.com/alireza-frj4/BlogBackEnd/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	database.Connect()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env file")
	}
	port := os.Getenv("PORT")
	app := fiber.New()
	routes.Setup(app)
	app.Listen(":" + port)
}
