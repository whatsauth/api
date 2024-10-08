package main

import (
	"log"

	"github.com/gofiber/fiber/v2/middleware/logger"

	"api/config"

	"api/helper/ws"

	"api/helper/chatroot"
	"api/helper/wa"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"api/url"

	"github.com/gofiber/fiber/v2"
)

func main() {
	/* 	clients, err := wa.ConnectAllClient(config.Mongoconn, config.ContainerDB)
	   	if err != nil {
	   		log.Panic(err)
	   	}

	   	config.MapClient.StoreAllClient(clients) */

	var err error
	config.Client, err = wa.ConnectAllClient(config.Mongoconn, config.ContainerDB)
	if err != nil {
		log.Panic(err)
	}

	go ws.RunHub()
	go chatroot.RunHub()

	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))

	site.Use(logger.New(logger.Config{
		Format: "${status} - ${method} ${path}\n",
	}))

	url.Web(site)
	log.Fatal(site.Listen(config.Port))
}
