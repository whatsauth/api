package controller

import (
	"api/config"
	"api/helper/wa"

	"github.com/aiteung/atdb"
	"github.com/aiteung/atmessage"
	"github.com/aiteung/musik"
	"github.com/whatsauth/watoken"
	"go.mau.fi/whatsmeow/types"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gofiber/fiber/v2"
)

type Header struct {
	Token string `reqHeader:"token"`
}

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

func SendMessage(c *fiber.Ctx) error {
	var h Header
	err := c.ReqHeaderParser(&h)
	if err != nil {
		return err
	}
	payload, err := watoken.Decode(config.PublicKey, h.Token)
	if err != nil {
		return err
	}
	_, err = atdb.GetOneLatestDoc[wa.User](config.Mongoconn, "user", bson.M{"phonenumber": payload.Id})
	var response atmessage.Response
	if err == nil {
		var notif atmessage.Notif
		err = c.BodyParser(&notif)
		if err != nil {
			return err
		}
		client := wa.GetWaClient(payload.Id, config.Client, config.Mongoconn)
		resp, _ := atmessage.SendMessage(notif.Messages, types.NewJID(notif.User, notif.Server), client.WAClient)

		response.Response = resp.ID
	}

	return c.JSON(response)
}
