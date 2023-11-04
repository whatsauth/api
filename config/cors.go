package config

import (
	"strings"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

var origins = []string{
	"https://wa.my.id",
	"https://whatsauth.my.id",
	"https://www.whatsauth.my.id",
}

var Cors = cors.Config{
	AllowOrigins:     strings.Join(origins[:], ","),
	AllowHeaders:     "Origin",
	ExposeHeaders:    "Content-Length",
	AllowCredentials: true,
}
