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
	if err == nil {
		var txt wa.TextMessage
		err = c.BodyParser(&txt)
		if err != nil {
			return err
		}
		client, _ := wa.GetWaClient(payload.Id, config.Client, config.Mongoconn, config.ContainerDB)
		resp, _ := wa.SendTextMessage(txt, client.WAClient)
		response.Response = resp.ID
	}

	return c.JSON(response)
}
