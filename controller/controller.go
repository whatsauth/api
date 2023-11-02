package controller

import (
	"api/helper"

	"github.com/aiteung/musik"
	"go.mau.fi/whatsmeow/types"

	"github.com/gofiber/fiber/v2"
)

func Homepage(c *fiber.Ctx) error {
	ipaddr := musik.GetIPaddress()
	return c.JSON(ipaddr)
}

func Device(c *fiber.Ctx) error {
	var phonejid = types.JID{
		User:   c.Params("+"),
		Server: "s.whatsapp.net",
	}
	qr := make(chan helper.QRStatus)

	go helper.GetQRString(helper.GetClient(phonejid), qr)
	a := <-qr
	return c.JSON(a)
}
