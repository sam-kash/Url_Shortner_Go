package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func setUpRoutes(app *fiber.App) { // This function has the list of all the routes
	app.Get("/:url", routes.ResolveURL)

	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	app := fiber.New()

	app.Use(Logger.New())

	setUpRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
