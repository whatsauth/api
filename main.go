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
	//privateKey, publicKey := watoken.GenerateKey()
	// loop database client untuk jalankan go helper.Connect(helper.GetClient(c.Params("+")), qr) disini pake for range
	wa.ConnectAllClient() // bahaya bug nya keluar semua device id nya
	//helper.Start(helper.GetClient("6287752000300"))
	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))
	url.Web(site)
	log.Fatal(site.Listen(config.Port))
}
