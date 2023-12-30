package controller

import (
	"api/config"
	"log"
	"net/url"

	"github.com/whatsauth/wa"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"

	"github.com/gofiber/fiber/v2"
)

func ClientList(c *fiber.Ctx) error {
	var resp wa.QRStatus
	payload, err := watoken.Decode(config.PublicKey, c.Params("+"))
	if err == nil {
		phonenumber := payload.Id
		qr := make(chan wa.QRStatus)

		waclient, _ := wa.GetWaClient(phonenumber, config.Client, config.Mongoconn, config.ContainerDB)
		//waclient, _ := wa.SetWaClient(phonenumber, config.Clients, config.Mongoconn, config.ContainerDB)
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
	//waclient, err := wa.GetWaClientMap(phonenumber, &config.MapClient, config.Mongoconn, config.ContainerDB)
	waclient, err := wa.GetWaClient(phonenumber, config.Client, config.Mongoconn, config.ContainerDB)
	//waclient, err := wa.SetWaClient(phonenumber, config.Clients, config.Mongoconn, config.ContainerDB)
	if err != nil {
		log.Println("error di getwaclientmap", err)
		resp = wa.QRStatus{Status: false, Message: err.Error()}
	} else {
		go wa.PairConnectStore(waclient, &config.MapClient, qr)
		resp = <-qr
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
		apdet := atdb.ReplaceOneDoc(config.Mongoconn, "user", bson.M{"phonenumber": payload.Id}, useraccount)
		if apdet.ModifiedCount == 0 {
			atdb.InsertOneDoc(config.Mongoconn, "user", useraccount)
		}

	}

	return c.JSON(useraccount)
}
