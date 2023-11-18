package controller

import (
	"api/config"

	"github.com/whatsauth/wa"

	"github.com/aiteung/atdb"
	"github.com/aiteung/atmessage"
	"github.com/gofiber/fiber/v2"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

func SendTextMessage(c *fiber.Ctx) error {
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
	response.Response = "WebHook Belum di daftarkan"
	if err == nil {
		var txt wa.TextMessage
		err = c.BodyParser(&txt)
		if err != nil {
			return err
		}
		var msg string
		if txt.Messages == "" {
			msg = "pesan kosong"
		} else {
			client, _ := wa.GetWaClient(payload.Id, config.Client, config.Mongoconn, config.ContainerDB)
			resp, _ := wa.SendTextMessage(txt, client.WAClient)

			if resp.Timestamp.IsZero() {
				msg = "device belum di start"
			} else {
				msg = "ID:" + resp.ID + " WARespon:" + resp.Timestamp.String() + " PeerTiming:" + resp.DebugTimings.PeerEncrypt.String() + " GetDeviceTiming:" + resp.DebugTimings.GetDevices.String()
			}
		}
		response.Response = msg
	}

	return c.JSON(response)
}
