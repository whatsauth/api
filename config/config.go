package config

import (
	"api/helper/wa"

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

var WhatsAuthPhoneNumber = "6283131895000"

var Port = ":8080"
