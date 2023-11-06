package controller

import (
	"api/config"
	"net/url"

	"github.com/whatsauth/wa"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/aiteung/atdb"
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

func SignUp(c *fiber.Ctx) error {
	var h Header
	var useraccount wa.User
	err := c.ReqHeaderParser(&h)
	if err != nil {
		return err
	}
	payload, err := watoken.Decode(config.PublicKey, h.Token)
	if err != nil {
		return err
	} else {
		var webhook wa.WebHook
		err = c.BodyParser(&webhook)
		if err != nil {
			return err
		}
		_, err := url.Parse(webhook.URL)
		if err != nil {
			return err
		}
		useraccount.PhoneNumber = payload.Id
		useraccount.WebHook = webhook
		newtoken, _ := watoken.EncodeforHours(payload.Id, config.PrivateKey, 720)
		useraccount.Token = newtoken
		olddata, err := atdb.GetOneLatestDoc[wa.User](config.Mongoconn, "user", bson.M{"phonenumber": payload.Id})
		if err != nil {
			atdb.InsertOneDoc(config.Mongoconn, "user", useraccount)
		} else {
			useraccount = olddata
		}

	}

	return c.JSON(useraccount)
}
