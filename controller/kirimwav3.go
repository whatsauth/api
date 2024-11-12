package controller

import (
	"api/config"
	"api/model"
	"log"

	"api/helper/wa"

	"api/helper/atdb"
	"api/helper/watoken"

	"github.com/gofiber/fiber/v2"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mongodb.org/mongo-driver/bson"
)

func SendDocumentMessageV3(c *fiber.Ctx) error {
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

func SendImageMessageV3(c *fiber.Ctx) error {
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

func SendTextMessageV3FromUser(c *fiber.Ctx) error {
	var h model.Header
	var resp whatsmeow.SendResponse

	err := c.ReqHeaderParser(&h)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userofficial, err := watoken.Decode(config.PublicKey, h.Token)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	var txt wa.TextMessage
	err = c.BodyParser(&txt)
	if err != nil {
		return err
	}
	//memastikan inputan nomor sesuai dengan format
	txt.To = formatPhoneNumber(txt.To)
	//ambil sender langganan dari log histori pengiriman sender official number
	// sender := log.GetSenderNumber(txt.To, config.Mongoconn)
	// ambil sender dari collection senderofficial number
	sender := wa.GetOfficialSenderNumber(userofficial.Id, config.Mongoconn)
	if sender == "" { //jika masih tidak ada maka user token tidak ada akses ke nomor official
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Token anda tidak memiliki akses ke nomor official"})
	}
	//kalo reset device nya ganti maka langsung aja acak
	if txt.Messages != "" {
		client, IsNewClient, err := wa.GetWaClient(sender, config.Client, config.Mongoconn, config.ContainerDB)
		if err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Waclient belum di start : " + err.Error()})
		}
		if IsNewClient {
			config.Client = append(config.Client, client)
		}
		//check untuk wa personal apakah nomornya sudah ada wa nya atau belum
		if !txt.IsGroup {
			onwa, err := client.WAClient.IsOnWhatsApp([]string{"+" + txt.To})
			if err != nil {
				//log error
				log.Println(userofficial)
				log.Println(sender)
				document := bson.D{
					{Key: "sender", Value: sender},
					{Key: "user", Value: userofficial.Id},
					{Key: "msg", Value: txt},
				}
				go atdb.InsertOneDoc(client.Mongoconn, "logerror", document)
				return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "nomor: " + sender + " linked Device belum diaktifkan : " + err.Error()})
			}
			if len(onwa) == 0 { //jika tidak terdeteksi sama sekali. nomor terlalu panjang
				return c.Status(fiber.StatusLengthRequired).JSON(fiber.Map{"error": "Nomor terlalu panjang"})
			}
			if !onwa[0].IsIn {
				return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Nomor tidak terdaftar di whatsapp"})
			}
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
			return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{"error": "Read wa tidak jalan : " + err.Error()})
			//return err
		}
		// simpen log siapa yang kirim siapa penerima nya
		go wa.LogSenderReceiverUpdate(sender, txt.To, config.Mongoconn)

		resp, err = wa.SendTextMessage(txt, client.WAClient)
		if err != nil {
			return c.Status(fiber.StatusFailedDependency).JSON(fiber.Map{"error": "Kirim pesan text wa ke websocket wa web tidak berhasil : " + err.Error()})
			//return err
		}
		//resp.Timestamp.IsZero()
	}

	return c.JSON(resp)
}
