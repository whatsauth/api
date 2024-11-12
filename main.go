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
	log.Println("Melakukan koneksi ke semua client yang sudah terdaftar di database")
	config.Client, err = wa.ConnectAllClient(config.Mongoconn, config.ContainerDB)
	if err != nil {
		log.Println("kesalahan ketika melakukan koneksi wa dari database")
		log.Panic(err)
	}

	log.Println("Menjalankan go routine untuk login qr")
	go ws.RunHub()
	log.Println("Menjalankan go routine untuk chatgpl")
	go chatroot.RunHub()

	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))

	site.Use(logger.New(logger.Config{
		Format: "${status} - ${method} ${path}\n",
	}))

	url.Web(site)
	log.Println("Menjalankan service fiber http pada: ", config.Port)
	log.Fatal(site.Listen(config.Port))
}
