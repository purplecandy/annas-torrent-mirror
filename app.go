package main

import (
	"annas-mirror/cache"
	"annas-mirror/database"
	"annas-mirror/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/pug/v2"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Println("Warning: Unable to load .env file")
	}

	cache.ConnectCache()
	database.ConnectDB()

	engine := pug.New("./views", ".pug")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", routes.AllTorrentGroups)
	app.Get("/:group", routes.TorrentGroupData)

	// action := actions.DispatchSyncTorrents()

	// cache.Dispatcher.Enqueue(action, asynq.Timeout(5*time.Minute))

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.Render("index", fiber.Map{"Title": "Hello World!"})
	// })

	app.Listen(":3000")
}
