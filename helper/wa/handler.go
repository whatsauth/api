package wa

import (
	"context"
	"fmt"
	"time"

	"api/model"

	"api/helper/atdb"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/proto"
)

func HandlingMessage(Info *types.MessageInfo, Message *waE2E.Message, client *WaClient) {
	if !Info.IsFromMe && (Info.Chat.Server != "broadcast") && (Info.Chat.User != "status") {
		var WAIface = model.IteungWhatsMeowConfig{
			Waclient: client.WAClient,
			Info:     Info,
			Message:  Message,
		}
		//simpan log pesan Message.ExtendedTextMessage.ContextInfo
		// go atdb.InsertOneDoc(client.Mongoconn, "inbox", WAIface.Message)
		//membuat struct untuk iteung v2
		Pesan := Whatsmeow2Struct(WAIface)
		//kirim ke webhook
		filter := bson.M{"phonenumber": client.PhoneNumber}
		userdt, _ := atdb.GetOneLatestDoc[User](client.Mongoconn, "user", filter)
		//go client.WAClient.SendChatPresence(Info.Chat, types.ChatPresenceComposing, types.ChatPresenceMediaText)
		go client.WAClient.MarkRead([]string{Info.ID}, time.Now(), Info.Chat, Info.Sender)
		result, err := PostStructWithToken[model.Response]("secret", userdt.WebHook.Secret, Pesan, userdt.WebHook.URL)
		if err != nil {
			var wamsg waE2E.Message
			wamsg.Conversation = proto.String(err.Error() + " RESULT:" + result.Response)
			client.WAClient.SendMessage(context.Background(), Info.Chat, &wamsg)
		}
	}
}

func Whatsmeow2Struct(WAIface model.IteungWhatsMeowConfig) (im model.IteungMessage) {
	im.Phone_number = GetPhoneNumber(WAIface)
	im.Chat_number = WAIface.Info.Chat.User
	im.Chat_server = WAIface.Info.Chat.Server
	im.Alias_name = WAIface.Info.PushName
	im.Message = GetMessage(WAIface.Message)
	im.EntryPoint = GetEntryPointDetail(WAIface)
	im.From_link = GetStatusFromLink(WAIface)
	if im.From_link {
		im.From_link_delay = GetFromLinkDelay(WAIface.Message)
	}
	im.Filename, im.Filedata = GetFile(WAIface.Waclient, WAIface.Message)
	im.Longitude, im.Latitude, im.LiveLoc = GetLongLat(WAIface.Message)
	if WAIface.Info.Chat.Server == "g.us" {
		groupInfo, err := WAIface.Waclient.GetGroupInfo(WAIface.Info.Chat)
		if err != nil {
			fmt.Println("cek err : ", err)
		}
		if groupInfo != nil {
			im.Group = groupInfo.GroupName.Name + "@" + WAIface.Info.Chat.User
			im.Group_name = groupInfo.GroupName.Name
			im.Group_id = WAIface.Info.Chat.User
		}
		im.Is_group = true
	}
	return
}
