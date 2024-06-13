package controller

import (
	"api/config"
	"api/model"

	"api/helper/wa"

	"api/helper/atdb"
	"api/helper/watoken"

	"github.com/gofiber/fiber/v2"
	"go.mau.fi/whatsmeow/types"
	"go.mongodb.org/mongo-driver/bson"
)

func SendDocumentMessage(c *fiber.Ctx) error {
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
	response.Response = "WebHook Belum di daftarkan"
	if err != nil {
		return c.JSON(response)
	}

	var doc wa.DocumentMessage
	err = c.BodyParser(&doc)
	if err != nil {
		return err
	}
	var msg string
	if doc.Base64Doc == "" {
		msg = "dokumen kosong"
	} else {
		client, IsNewClient, err := wa.GetWaClient(payload.Id, config.Client, config.Mongoconn, config.ContainerDB)
		if err != nil {
			return err
		}
		if IsNewClient {
			config.Client = append(config.Client, client)
		}

		var targetjid types.JID
		targetjid.User = doc.To
		if doc.IsGroup {
			targetjid.Server = "g.us"
		} else {
			targetjid.Server = "s.whatsapp.net"
		}
		err = client.WAClient.SendChatPresence(targetjid, types.ChatPresenceComposing, types.ChatPresenceMediaText)
		if err != nil {
			return err
		}

		resp, err := wa.SendDocumentMessage(doc, client.WAClient)
		if err != nil {
			return err
		}
		if resp.Timestamp.IsZero() {
			msg = "device belum di start"
		} else {
			msg = "ID:" + resp.ID + " WARespon:" + resp.Timestamp.String() + " PeerTiming:" + resp.DebugTimings.PeerEncrypt.String() + " GetDeviceTiming:" + resp.DebugTimings.GetDevices.String()
		}
	}
	response.Response = msg

	return c.JSON(response)
}

func SendImageMessage(c *fiber.Ctx) error {
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
	response.Response = "WebHook Belum di daftarkan"
	if err != nil {
		return c.JSON(response)
	}

	var img wa.ImageMessage
	err = c.BodyParser(&img)
	if err != nil {
		return err
	}
	var msg string
	if img.Base64Image == "" {
		msg = "gambar kosong"
	} else {
		client, IsNewClient, err := wa.GetWaClient(payload.Id, config.Client, config.Mongoconn, config.ContainerDB)
		if err != nil {
			return err
		}
		if IsNewClient {
			config.Client = append(config.Client, client)
		}

		var targetjid types.JID
		targetjid.User = img.To
		if img.IsGroup {
			targetjid.Server = "g.us"
		} else {
			targetjid.Server = "s.whatsapp.net"
		}
		err = client.WAClient.SendChatPresence(targetjid, types.ChatPresenceComposing, types.ChatPresenceMediaText)
		if err != nil {
			return err
		}

		resp, err := wa.SendImageMessage(img, client.WAClient)
		if err != nil {
			return err
		}
		if resp.Timestamp.IsZero() {
			msg = "device belum di start"
		} else {
			msg = "ID:" + resp.ID + " WARespon:" + resp.Timestamp.String() + " PeerTiming:" + resp.DebugTimings.PeerEncrypt.String() + " GetDeviceTiming:" + resp.DebugTimings.GetDevices.String()
		}
	}
	response.Response = msg

	return c.JSON(response)
}

func SendTextMessage(c *fiber.Ctx) error {
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
	response.Response = "WebHook Belum di daftarkan"
	if err != nil {
		return c.JSON(response)
	}

	var txt wa.TextMessage
	err = c.BodyParser(&txt)
	if err != nil {
		return err
	}
	var msg string
	if txt.Messages == "" {
		msg = "pesan kosong"
	} else {
		client, IsNewClient, err := wa.GetWaClient(payload.Id, config.Client, config.Mongoconn, config.ContainerDB)
		if err != nil {
			return err
		}
		if IsNewClient {
			config.Client = append(config.Client, client)
		}

		var targetjid types.JID
		targetjid.User = txt.To
		if txt.IsGroup {
			targetjid.Server = "g.us"
		} else {
			targetjid.Server = "s.whatsapp.net"
		}
		err = client.WAClient.SendChatPresence(targetjid, types.ChatPresenceComposing, types.ChatPresenceMediaText)
		if err != nil {
			return err
		}

		resp, err := wa.SendTextMessage(txt, client.WAClient)
		if err != nil {
			return err
		}
		if resp.Timestamp.IsZero() {
			msg = "device belum di start"
		} else {
			msg = "ID:" + resp.ID + " WARespon:" + resp.Timestamp.String() + " PeerTiming:" + resp.DebugTimings.PeerEncrypt.String() + " GetDeviceTiming:" + resp.DebugTimings.GetDevices.String()
		}
	}
	response.Response = msg

	return c.JSON(response)
}
