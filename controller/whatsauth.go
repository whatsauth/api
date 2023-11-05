package controller

import (
	"api/config"
	"api/helper/wa"
	"api/helper/ws"
	"fmt"

	"github.com/aiteung/atdb"
	"github.com/aiteung/atmessage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

func PostWhatsAuthRequest(c *fiber.Ctx) error {
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
		var req ws.WhatsauthRequest
		err := c.BodyParser(&req)
		if err != nil {
			return err
		}
		login, err := watoken.Encode(req.Phonenumber, config.PrivateKey)
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
			Messages: "Berhasil",
		}
		wsstatus := ws.SendStructTo(req.Uuid, infologin)
		if !wsstatus {
			txt.Messages = "Status Socket Tertutup"
		}
		if req.Uuid[0:1] == "m" {
			txt.Messages = "\n Selanjutnya kakak *tinggal login* ya kak."
			tokenstring, er := watoken.EncodeforSeconds(req.Phonenumber, config.PrivateKey, 30)
			if er != nil {
				return er
			}

			urlakses := watoken.GetAppUrl(req.Uuid) + "?uuid=" + tokenstring
			txt.Messages += fmt.Sprintf("Magic Link : %v", urlakses)
		}
		client := wa.GetWaClient(payload.Id, config.Client, config.Mongoconn)
		resp, _ := wa.SendTextMessage(txt, client.WAClient)
		response.Response = resp.ID

	}

	return c.JSON(response)
}

func WsWhatsAuth(c *websocket.Conn) {
	ws.RunSocket(c, config.PublicKey, config.PrivateKey)
}
