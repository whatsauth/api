package config

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
)

/* var origins = []string{
	"https://wa.my.id",
	"https://whatsauth.my.id",
	"https://www.whatsauth.my.id",
	"https://my.wa.my.id",
	"https://lapor.acad-csirt.org/",
} */

var Cors = cors.Config{
	AllowOrigins:     "*", //strings.Join(origins[:], ","),
	AllowHeaders:     "Origin, Token, Content-Type",
	ExposeHeaders:    "Content-Length",
	AllowCredentials: true,
}
