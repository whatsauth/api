package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"api/config"
	"api/helper/chatroot"
	"api/helper/wa"
	"api/helper/ws"
	"api/url"
)

func main() {
	// Konfigurasi logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	log.Info().Msg("Aplikasi sedang berjalan...")

	// Cek koneksi ke database
	log.Info().Msg("Melakukan koneksi ke semua client yang sudah terdaftar di database")
	var err error
	config.Client, err = wa.ConnectAllClient(config.Mongoconn, config.ContainerDB)
	if err != nil {
		log.Fatal().Err(err).Msg("Kesalahan ketika melakukan koneksi WA dari database")
	}
	// Menjalankan Go Routine
	log.Info().Msg("Menjalankan go routine untuk login QR")
	go ws.RunHub()

	log.Info().Msg("Menjalankan go routine untuk ChatGPT")
	go chatroot.RunHub()

	// Setup Fiber
	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))

	// Middleware logging request
	site.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))

	url.Web(site)
	// Jalankan server
	log.Info().Msgf("Menjalankan service fiber HTTP pada: %s", config.Port)
	err = site.Listen(config.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("Server gagal berjalan")
	}
}
