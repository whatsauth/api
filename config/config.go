package config

import (
	"api/helper/wa"

	"github.com/gofiber/fiber/v2"
	"go.mau.fi/whatsmeow/types"
)

var Client []*wa.WaClient

var Iteung = fiber.Config{
	Prefork:       false,
	CaseSensitive: true,
	StrictRouting: true,
	ServerHeader:  "WhatsAuth",
	AppName:       "API Message Router",
}

var PhoneNumber = types.JID{
	User:   "6287752000300",
	Server: "s.whatsapp.net",
}

var Port = ":8080"
