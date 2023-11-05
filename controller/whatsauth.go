package controller

import (
	"api/config"
	"api/helper/wa"
	"api/helper/ws"

	"github.com/aiteung/atmessage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/whatsauth/watoken"
)

func PostWhatsAuthRequest(c *fiber.Ctx) error {
	var h Header
	err := c.ReqHeaderParser(&h)
	if err != nil {
		return err
	}
	var req ws.WhatsauthRequest
	var response atmessage.Response
	if h.Token == config.WhatsAuthReqToken {
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
		ws.SendStructTo(req.Uuid, infologin)
		client := wa.GetWaClient(config.WhatsAuthPhoneNumber, config.Client, config.Mongoconn)
		resp, _ := wa.SendTextMessage(txt, client.WAClient)
		response.Response = resp.ID

	}

	return c.JSON(response)
}

func WsWhatsAuth(c *websocket.Conn) {
	ws.RunSocket(c, config.PublicKey)
}
