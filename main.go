package main

import (
	"log"

	"api/config"
	"api/helper"

	"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"api/url"

	"github.com/gofiber/fiber/v2"
)

func main() {
	helper.GetQRString(helper.GetClient())
	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))
	url.Web(site)
	log.Fatal(site.Listen(musik.Dangdut()))
}
