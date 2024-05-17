package wa

import (
	"context"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func SendTextMessage(txt TextMessage, whatsapp *whatsmeow.Client) (resp whatsmeow.SendResponse, err error) {
	server := "s.whatsapp.net"
	if txt.IsGroup {
		server = "g.us"
	}
	//go whatsapp.SendChatPresence(types.NewJID(txt.To, server), types.ChatPresenceComposing, types.ChatPresenceMediaText)
	var wamsg waProto.Message
	wamsg.Conversation = proto.String(txt.Messages)
	resp, err = whatsapp.SendMessage(context.Background(), types.NewJID(txt.To, server), &wamsg)
	return resp, err
}
