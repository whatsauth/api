package wa

import (
	"context"
	"time"

	"github.com/aiteung/atdb"
	"github.com/aiteung/atmessage"
	"github.com/aiteung/module"
	"github.com/aiteung/module/model"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/proto"
)

func HandlingMessage(Info *types.MessageInfo, Message *waProto.Message, client *WaClient) {
	go client.WAClient.MarkRead([]string{Info.ID}, time.Now(), Info.Chat, Info.Sender)
	if !Info.IsFromMe && (Info.Chat.Server != "broadcast") && (Info.Chat.User != "status") {
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
		go client.WAClient.SendChatPresence(Info.Chat, types.ChatPresenceComposing, types.ChatPresenceMediaText)
		result, err := PostStructWithToken[atmessage.Response]("secret", userdt.WebHook.Secret, Pesan, userdt.WebHook.URL)
		if err != nil {
			var wamsg waProto.Message
			wamsg.Conversation = proto.String(err.Error() + " RESULT:" + result.Response)
			client.WAClient.SendMessage(context.Background(), Info.Chat, &wamsg)
		}
	}
}
