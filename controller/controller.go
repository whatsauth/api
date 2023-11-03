package controller

import (
	"api/config"
	"api/helper"

	"github.com/aiteung/musik"
	"github.com/whatsauth/watoken"

	"github.com/gofiber/fiber/v2"
)

func Homepage(c *fiber.Ctx) error {
	ipaddr := musik.GetIPaddress()
	return c.JSON(ipaddr)
}

func Device(c *fiber.Ctx) error {
	body, err := watoken.Decode(config.PublicKey, c.Params("+"))
	if err != nil {
		return err
	}
	qr := make(chan helper.QRStatus)

	go helper.Connect(body.Id, qr)
	a := <-qr
	return c.JSON(a)
}
