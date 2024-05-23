package controller

import (
	"api/config"
	"api/helper"
	"api/model"

	"api/helper/ws"

	"api/helper/wa"

	"api/helper/watoken"

	"api/helper/atdb"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func PostWhatsAuthRequest(c *fiber.Ctx) error {
	var h model.Header
	err := c.ReqHeaderParser(&h)
	if err != nil {
		return err
	}
	payload, err := watoken.Decode(config.PublicKey, h.Token)
	if err != nil {
		return err
	}
	_, err = atdb.GetOneLatestDoc[wa.User](config.Mongoconn, "user", bson.M{"phonenumber": payload.Id})
	var response model.Response
	if err == nil {
		var req ws.WhatsauthRequest
		err := c.BodyParser(&req)
		if err != nil {
			return err
		}
		login, err := watoken.EncodeforHours(req.Phonenumber, req.Aliasname, config.PrivateKey, 18)
		if err != nil {
			return err
		}
		infologin := ws.LoginInfo{
			Phone: req.Phonenumber,
			Login: login,
			Uuid:  req.Uuid,
		}
		var txt = wa.TextMessage{
			To:       req.Phonenumber,
			IsGroup:  false,
			Messages: "*WhatsAuth Indonesia Login*\n",
		}
		wsstatus := ws.SendStructTo(req.Uuid, infologin)
		//go atdb.InsertOneDoc(config.Mongoconn, "logwauth", infologin)
		if !wsstatus && req.Uuid[0:1] != "m" {
			txt.Messages += "Sesi QR sudah habis, mohon pastikan memiliki waktu cukup untuk scan QR."
		} else if req.Uuid[0:1] == "m" {
			txt.Messages += "Selanjutnya kakak klik saja magic link di bawah ini ya kak:\n"
			tokenstring, er := watoken.EncodeforSeconds(req.Phonenumber, req.Aliasname, config.PrivateKey, 30)
			if er != nil {
				return er
			}

			urlakses := watoken.GetAppUrl(req.Uuid) + "?uuid=" + tokenstring
			txt.Messages += urlakses
		} else {
			filter := bson.M{"notif": "login"}
			pick := helper.PickPantun(filter)
			txt.Messages += pick + "Yey... login diterima kak, silahkan kembali ke browser lagi ya."
		}
		client, IsNewClient, err := wa.GetWaClient(payload.Id, config.Client, config.Mongoconn, config.ContainerDB)
		if err != nil {
			return err
		}
		if IsNewClient {
			config.Client = append(config.Client, client)
		}

		resp, _ := wa.SendTextMessage(txt, client.WAClient)
		response.Response = resp.ID

	}

	return c.JSON(response)
}

func WsWhatsAuth(c *websocket.Conn) {
	ws.RunSocket(c, config.PublicKey, config.PrivateKey)
}
