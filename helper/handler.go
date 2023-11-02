package helper

import (
	"time"

	"github.com/aiteung/atdb"
	"github.com/aiteung/module"
	"github.com/aiteung/module/model"
	"github.com/google/uuid"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func HandlingMessage(Info *types.MessageInfo, Message *waProto.Message, Client *whatsmeow.Client) {
	go Client.MarkRead([]string{Info.ID}, time.Now(), Info.Chat, Info.Sender)
	if !Info.IsFromMe {
		var WAIface = model.IteungWhatsMeowConfig{
			Waclient: Client,
			Info:     Info,
			Message:  Message,
		}
		//membuat struct untuk iteung v2
		Pesan := module.Whatsmeow2Struct(WAIface)
		if Pesan.Phone_number == "6285156007137" || Pesan.Phone_number == "6285757707248" || Pesan.Phone_number == "6281312000300" {
			go atdb.InsertOneDoc(config.MongoIteungConn, "log_iteung_message_gowa", ITeungMessageUUID{
				UUID:          uuid.New().String(),
				IteungMessage: Pesan,
			})
		}
		//kirim ke backend iteung v2
		resp, errmessage := module.SendToIteungAPI(Pesan, config.IteungURLV2)
		//log error untuk debug
		if errmessage != "" {
			go atdb.InsertOneDoc(config.MongoIteungConn, "log_error", ITeungErrorLog{
				UUID:         uuid.New().String(),
				Container:    "iteung-gowa",
				Message:      Pesan,
				Response:     resp,
				ErrorMessage: errmessage,
			})
		}
		if resp.Response != "" {
			go Client.SendChatPresence(Info.Chat, "composing", "")
		}
		//fmt.Println("respon backend : ", resp)
	}
}
