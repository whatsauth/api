package wa

import (
	"time"

	"github.com/aiteung/atapi"
	"github.com/aiteung/atdb"
	"github.com/aiteung/atmessage"
	"github.com/aiteung/module"
	"github.com/aiteung/module/model"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mongodb.org/mongo-driver/bson"
)

func HandlingMessage(Info *types.MessageInfo, Message *waProto.Message, client *WaClient) {
	go client.WAClient.MarkRead([]string{Info.ID}, time.Now(), Info.Chat, Info.Sender)
	if !Info.IsFromMe {
		var WAIface = model.IteungWhatsMeowConfig{
			Waclient: client.WAClient,
			Info:     Info,
			Message:  Message,
		}
		//membuat struct untuk iteung v2
		Pesan := module.Whatsmeow2Struct(WAIface)
		//kirim ke webhook
		filter := bson.M{"phonenumber": client.PhoneNumber}
		userdt, _ := atdb.GetOneLatestDoc[User](client.Mongoconn, "user", filter)
		atapi.PostStructWithToken[atmessage.Response]("secret", userdt.WebHook.Secret, Pesan, userdt.WebHook.URL)
	}
}
