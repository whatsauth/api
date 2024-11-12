package config

import (
	"api/helper/wa"

	"github.com/gofiber/fiber/v2"
)

var Client []*wa.WaClient

var MapClient = wa.NewMapClient()

//var Clients *wa.Clients

var Iteung = fiber.Config{
	Prefork:       false,
	CaseSensitive: true,
	StrictRouting: true,
	ServerHeader:  "WhatsAuth",
	AppName:       "API Message Router",
}

var Port = "127.0.0.1:8080"
