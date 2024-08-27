package main

import (
	"annas-mirror/actions"
	"annas-mirror/cache"
	"annas-mirror/database"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Println("Warning: Unable to load .env file")
	}

	cache.ConnectCache()
	database.ConnectDB()

	app := fiber.New()

	action := actions.DispatchSyncTorrents()

	cache.Dispatcher.Enqueue(action)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}
