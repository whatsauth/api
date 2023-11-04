package controller

import (
	"api/config"
	"api/helper/wa"

	"github.com/aiteung/musik"
	"github.com/whatsauth/watoken"

	"github.com/gofiber/fiber/v2"
)

func Homepage(c *fiber.Ctx) error {
	ipaddr := musik.GetIPaddress()
	return c.JSON(ipaddr)
}

func Device(c *fiber.Ctx) error {
	phonenumber := watoken.DecodeGetId(config.PublicKey, c.Params("+"))
	qr := make(chan wa.QRStatus)

	waclient := wa.GetWaClient(phonenumber, config.Client, config.Mongoconn)
	//go wa.QRConnect(waclient, qr)
	go wa.PairConnect(waclient, qr)
	a := <-qr
	return c.JSON(a)
}
