package config

import (
	"github.com/whatsauth/wa"

	"github.com/gofiber/fiber/v2"
)

var Client []*wa.WaClient

var Iteung = fiber.Config{
	Prefork:       false,
	CaseSensitive: true,
	StrictRouting: true,
	ServerHeader:  "WhatsAuth",
	AppName:       "API Message Router",
}

var Port = ":8080"
