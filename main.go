package main

import (
	"log"

	"api/config"
	"api/helper/wa"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"api/url"

	"github.com/gofiber/fiber/v2"
)

func main() {
	wa.ConnectAllClient(config.Mongoconn)

	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))
	url.Web(site)
	log.Fatal(site.Listen(config.Port))
}
