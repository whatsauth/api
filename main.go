package main

import (
	"log"

	"api/config"

	"github.com/whatsauth/ws"

	"github.com/whatsauth/wa"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"api/url"

	"github.com/gofiber/fiber/v2"
)

func main() {
	var err error
	config.Client, err = wa.ConnectAllClient(config.Mongoconn, config.ContainerDB)
	if err != nil {
		log.Panic(err)
	}
	go ws.RunHub()

	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))
	url.Web(site)
	log.Fatal(site.Listen(config.Port))
}
