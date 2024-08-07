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

func CheckNumbersInWhatsApp(c *fiber.Ctx) error {
	var h model.Header

	err := c.ReqHeaderParser(&h)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	payload, err := watoken.Decode(config.PublicKey, h.Token)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	var phlist wa.PhoneList
	err = c.BodyParser(&phlist)
	if err != nil {
		return err
	}

	client, IsNewClient, err := wa.GetWaClient(payload.Id, config.Client, config.Mongoconn, config.ContainerDB)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Waclient belum di start : " + err.Error()})
	}
	if IsNewClient {
		config.Client = append(config.Client, client)
	}
	//memastikan inputan nomor sesuai dengan format
	//txt.To = formatPhoneNumber(txt.To)
	//check untuk wa personal apakah nomornya sudah ada wa nya atau belum
	var walist model.WaList

	walist.PhoneNumbers, err = client.WAClient.IsOnWhatsApp(phlist.PhoneNumbers)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "IsOnWhatsApp: " + err.Error()})
	}

	return c.JSON(walist)
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
