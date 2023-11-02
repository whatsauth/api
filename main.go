package main

import (
	"log"

	"api/config"

	"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"api/url"

	"github.com/gofiber/fiber/v2"
)

func main() {
	//privateKey, publicKey := watoken.GenerateKey()
	// loop database client untuk jalankan go helper.Connect(helper.GetClient(c.Params("+")), qr) disini pake for range
	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))
	url.Web(site)
	log.Fatal(site.Listen(musik.Dangdut()))
}
