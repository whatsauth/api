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
	var resp wa.QRStatus
	payload, err := watoken.Decode(config.PublicKey, c.Params("+"))
	if err == nil {
		phonenumber := payload.Id
		qr := make(chan wa.QRStatus)

		waclient := wa.GetWaClient(phonenumber, config.Client, config.Mongoconn)
		//go wa.QRConnect(waclient, qr)
		go wa.PairConnect(waclient, qr)
		resp = <-qr

	} else {
		resp = wa.QRStatus{Status: false, Message: "tidak terdaftar"}
	}

	return c.JSON(resp)
}
