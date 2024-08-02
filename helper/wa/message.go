package wa

import (
	"api/model"
	"context"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func SendDocumentMessage(doc DocumentMessage, whatsapp *whatsmeow.Client) (resp whatsmeow.SendResponse, err error) {
	server := "s.whatsapp.net"
	if doc.IsGroup {
		server = "g.us"
	}

	docData, err := base64.StdEncoding.DecodeString(doc.Base64Doc)
	if err != nil {
		return resp, err
	}

	respupload, err := whatsapp.Upload(context.Background(), docData, whatsmeow.MediaDocument)
	if err != nil {
		return resp, err
	}

	docMsg := &waE2E.DocumentMessage{
		Caption:       proto.String(doc.Caption),
		Mimetype:      proto.String(http.DetectContentType(docData)),
		FileName:      &doc.Filename,
		URL:           &respupload.URL,
		DirectPath:    &respupload.DirectPath,
		MediaKey:      respupload.MediaKey,
		FileEncSHA256: respupload.FileEncSHA256,
		FileSHA256:    respupload.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(docData))),
	}
	docMessage := &waE2E.Message{
		DocumentMessage: docMsg,
	}
	resp, err = whatsapp.SendMessage(context.Background(), types.NewJID(doc.To, server), docMessage)
	return resp, err
}

func SendImageMessage(img ImageMessage, whatsapp *whatsmeow.Client) (resp whatsmeow.SendResponse, err error) {
	server := "s.whatsapp.net"
	if img.IsGroup {
		server = "g.us"
	}

	imageData, err := base64.StdEncoding.DecodeString(img.Base64Image)
	if err != nil {
		return resp, err
	}

	respupload, err := whatsapp.Upload(context.Background(), imageData, whatsmeow.MediaImage)
	if err != nil {
		return resp, err
	}

	imgMsg := &waE2E.ImageMessage{
		Caption:       proto.String(img.Caption),
		URL:           proto.String(respupload.URL),
		DirectPath:    proto.String(respupload.DirectPath),
		MediaKey:      respupload.MediaKey,
		Mimetype:      proto.String(http.DetectContentType(imageData)),
		FileEncSHA256: respupload.FileEncSHA256,
		FileSHA256:    respupload.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(imageData))),
	}

	imgMessage := &waE2E.Message{
		ImageMessage: imgMsg,
	}
	resp, err = whatsapp.SendMessage(context.Background(), types.NewJID(img.To, server), imgMessage)
	return resp, err
}

func SendTextMessage(txt TextMessage, whatsapp *whatsmeow.Client) (resp whatsmeow.SendResponse, err error) {
	server := "s.whatsapp.net"
	if txt.IsGroup {
		server = "g.us"
	}
	//go whatsapp.SendChatPresence(types.NewJID(txt.To, server), types.ChatPresenceComposing, types.ChatPresenceMediaText)
	var wamsg waE2E.Message
	wamsg.Conversation = proto.String(txt.Messages)
	resp, err = whatsapp.SendMessage(context.Background(), types.NewJID(txt.To, server), &wamsg)
	return resp, err
}

func GetPhoneNumber(WAIface model.IteungWhatsMeowConfig) (phonenumber string) {
	phonenumber = WAIface.Info.Sender.User
	if WAIface.Message.ExtendedTextMessage != nil {
		if WAIface.Message.ExtendedTextMessage.ContextInfo != nil {
			//Kalo pake whatsapp Desktop masuk module ExtendedTextMessage ContextInfo expiration:0
			if WAIface.Message.ExtendedTextMessage.ContextInfo.Participant != nil {
				phonenumber = strings.Split(*WAIface.Message.ExtendedTextMessage.ContextInfo.Participant, "@")[0]
			}
		}
	}

	return

}

func GetMessage(Message *waE2E.Message) (message string) {
	switch {
	case Message.ExtendedTextMessage != nil:
		message = *Message.ExtendedTextMessage.Text
	case Message.DocumentMessage != nil:
		if Message.DocumentMessage.Caption != nil {
			message = *Message.DocumentMessage.Caption
		}
	case Message.ImageMessage != nil:
		if Message.ImageMessage.Caption != nil {
			message = *Message.ImageMessage.Caption
		}
	case Message.LiveLocationMessage != nil:
		message = Message.LiveLocationMessage.GetCaption()
	default:
		message = Message.GetConversation()
	}
	return

}

func GetLongLat(Message *waE2E.Message) (long, lat float64, liveloc bool) {
	if Message.ExtendedTextMessage != nil {
		if Message.ExtendedTextMessage.ContextInfo != nil {
			if Message.ExtendedTextMessage.ContextInfo.Participant != nil {
				if Message.ExtendedTextMessage.ContextInfo.QuotedMessage.LiveLocationMessage != nil {
					lat = *Message.ExtendedTextMessage.ContextInfo.QuotedMessage.LiveLocationMessage.DegreesLatitude
					long = *Message.ExtendedTextMessage.ContextInfo.QuotedMessage.LiveLocationMessage.DegreesLongitude
					liveloc = true
				}
			}

		}
	} else if Message.LiveLocationMessage != nil {
		long = *Message.LiveLocationMessage.DegreesLongitude
		lat = *Message.LiveLocationMessage.DegreesLatitude
		liveloc = true
	} else if Message.LocationMessage != nil {
		long = *Message.LocationMessage.DegreesLongitude
		lat = *Message.LocationMessage.DegreesLatitude
	}
	return
}

func GetFile(client *whatsmeow.Client, Message *waE2E.Message) (filename, filedata string) {
	if extMsg := Message.GetExtendedTextMessage(); extMsg != nil {
		if extMsg.ContextInfo == nil {
			return
		}
		if extMsg.ContextInfo.Participant == nil {
			return
		}
		if extMsg.ContextInfo.QuotedMessage.DocumentMessage != nil {
			filename = *extMsg.ContextInfo.QuotedMessage.DocumentMessage.DirectPath
			payload, err := client.Download(extMsg.ContextInfo.QuotedMessage.DocumentMessage)
			if err != nil {
				return
			}
			filedata = base64.StdEncoding.EncodeToString(payload)
		}
		if extMsg.ContextInfo.QuotedMessage.DocumentWithCaptionMessage != nil {
			filename = *extMsg.ContextInfo.QuotedMessage.DocumentWithCaptionMessage.Message.DocumentMessage.DirectPath
			payload, err := client.Download(extMsg.ContextInfo.QuotedMessage.DocumentWithCaptionMessage.Message.DocumentMessage)
			if err != nil {
				return
			}
			filedata = base64.StdEncoding.EncodeToString(payload)
		}
	} else if doc := Message.GetDocumentMessage(); doc != nil {
		switch {
		case doc.Title != nil:
			filename = *doc.Title
		case doc.FileName != nil:
			filename = *doc.FileName
		}
		payload, err := client.Download(doc)
		if err != nil {
			return
		}
		filedata = base64.StdEncoding.EncodeToString(payload)
	} else if img := Message.GetImageMessage(); img != nil {
		filename = strings.ReplaceAll(*img.Mimetype, "/", ".")
		filedata = GetBase64Filedata(img.URL, img.GetMediaKey())
		payload, err := client.Download(img)
		if err != nil {
			return
		}
		filedata = base64.StdEncoding.EncodeToString(payload)
	} else if docCap := Message.GetDocumentWithCaptionMessage(); docCap != nil {
		if docCap.GetMessage() == nil {
			return
		}
		switch {
		case docCap.GetMessage().GetDocumentMessage().Title != nil:
			filename = docCap.GetMessage().GetDocumentMessage().GetTitle()
		case docCap.GetMessage().GetDocumentMessage().FileName != nil:
			filename = docCap.GetMessage().GetDocumentMessage().GetFileName()
		}
		payload, err := client.Download(docCap.Message.DocumentMessage)
		if err != nil {
			return
		}

		filedata = base64.StdEncoding.EncodeToString(payload)
	}
	return

}

func GetStatusFromLink(WAIface model.IteungWhatsMeowConfig) (whmsg bool) {
	if WAIface.Message.ExtendedTextMessage != nil && WAIface.Info.Chat.Server == "s.whatsapp.net" {
		if WAIface.Message.ExtendedTextMessage.ContextInfo != nil {
			if WAIface.Message.ExtendedTextMessage.ContextInfo.EntryPointConversionSource != nil {
				msg := *WAIface.Message.ExtendedTextMessage.ContextInfo.EntryPointConversionSource
				if msg == "click_to_chat_link" {
					whmsg = true
				}
			}
		}
	}
	return
}

func GetEntryPointDetail(WAIface model.IteungWhatsMeowConfig) (details string) {
	if WAIface.Message.GetExtendedTextMessage() != nil {
		if WAIface.Message.GetExtendedTextMessage().GetContextInfo() != nil {
			sumber := WAIface.Message.GetExtendedTextMessage().GetContextInfo().GetEntryPointConversionSource()
			app := WAIface.Message.GetExtendedTextMessage().GetContextInfo().GetEntryPointConversionApp()
			delay := WAIface.Message.GetExtendedTextMessage().GetContextInfo().GetEntryPointConversionDelaySeconds()
			if sumber != "" {
				details = sumber + "|" + app + "|" + strconv.FormatUint(uint64(delay), 10)
			}
		}
	}
	return
}

func GetFromLinkDelay(Message *waE2E.Message) uint32 {
	return *Message.ExtendedTextMessage.ContextInfo.EntryPointConversionDelaySeconds
}
