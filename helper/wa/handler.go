package wa

import (
	"fmt"
	"time"

	"github.com/aiteung/atapi"
	"github.com/aiteung/atmessage"
	"github.com/aiteung/module"
	"github.com/aiteung/module/model"
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
		//kirim ke webhook
		resp, err := atapi.PostStructWithToken[atmessage.Response]("secret", "okokok", Pesan, "https://eov6tgpfbhsve67.m.pipedream.net")
		if err != "" {
			fmt.Println(err)
		}
		if resp.Response != "" {
			go Client.SendChatPresence(Info.Chat, "composing", "")
		}
	}
}
