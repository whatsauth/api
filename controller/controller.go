package controller

import (
	"api/helper"

	"github.com/aiteung/musik"

	"github.com/gofiber/fiber/v2"
)

func Homepage(c *fiber.Ctx) error {
	ipaddr := musik.GetIPaddress()
	return c.JSON(ipaddr)
}

func Device(c *fiber.Ctx) error {
	qr := make(chan helper.QRStatus)

	go helper.GetQRString(helper.GetClient(c.Params("+")), qr)
	a := <-qr
	return c.JSON(a)
}
