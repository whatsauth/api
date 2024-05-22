package controller

import (
	"api/config"
	"api/model"
	"net/url"

	"api/helper/wa"

	"go.mongodb.org/mongo-driver/bson"

	"api/helper/atdb"
	"api/helper/watoken"

	"github.com/gofiber/fiber/v2"
)

func ClientList(c *fiber.Ctx) error {
	var resp wa.QRStatus
	payload, err := watoken.Decode(config.PublicKey, c.Params("+"))
	if err == nil {
		phonenumber := payload.Id
		qr := make(chan wa.QRStatus)

		waclient, IsNewClient, err := wa.GetWaClient(phonenumber, config.Client, config.Mongoconn, config.ContainerDB)
		if err != nil {
			resp = wa.QRStatus{Status: false, Message: "nomor tidak terdaftar"}
			return c.JSON(resp)
		}
		if IsNewClient {
			config.Client = append(config.Client, waclient)
		}
		go wa.PairConnect(waclient, qr)
		resp = <-qr

	} else {
		resp = wa.QRStatus{Status: false, Message: "nomor tidak terdaftar"}
	}

	return c.JSON(resp)
}

func StartDevice(c *fiber.Ctx) error {
	var resp wa.QRStatus
	payload, err := watoken.Decode(config.PublicKey, c.Params("+"))
	if err == nil {
		phonenumber := payload.Id
		config.Client, err = wa.ConnectAllClient(config.Mongoconn, config.ContainerDB)
		if err != nil {
			resp = wa.QRStatus{Status: false, Message: "ResetDevice:" + err.Error()}
		} else {
			resp = wa.QRStatus{PhoneNumber: phonenumber, Status: false, Message: "Reset Device selesai"}
		}

	} else {
		resp = wa.QRStatus{PhoneNumber: c.Params("+"), Code: err.Error(), Status: false, Message: "nomor tidak terdaftar"}
	}

	return c.JSON(resp)
}

func Device(c *fiber.Ctx) error {
	var resp wa.QRStatus

	payload, err := watoken.Decode(config.PublicKey, c.Params("+"))
	if err != nil {
		resp = wa.QRStatus{Status: false, Message: "nomor tidak terdaftar"}
		return c.JSON(resp)
	}

	phonenumber := payload.Id
	qr := make(chan wa.QRStatus)
	waclient, IsNewClient, err := wa.GetWaClient(phonenumber, config.Client, config.Mongoconn, config.ContainerDB)
	if err != nil {
		resp = wa.QRStatus{Status: false, Message: err.Error()}
		return c.JSON(resp)
	}
	if IsNewClient {
		config.Client = append(config.Client, waclient)
	}
	go wa.PairConnectStore(waclient, &config.MapClient, qr)
	resp = <-qr

	return c.JSON(resp)
}

func SignUp(c *fiber.Ctx) error {
	var h model.Header
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
		newtoken, _ := watoken.EncodeforHours(payload.Id, payload.Alias, config.PrivateKey, 720)
		useraccount.Token = newtoken
		apdet, _ := atdb.ReplaceOneDoc(config.Mongoconn, "user", bson.M{"phonenumber": payload.Id}, useraccount)
		if apdet.ModifiedCount == 0 {
			atdb.InsertOneDoc(config.Mongoconn, "user", useraccount)
		}

	}

	return c.JSON(useraccount)
}
