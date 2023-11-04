package ws

import (
	"github.com/whatsauth/watoken"
)

func MagicLinkEvent(roomId string, PublicKey string) {
	if roomId[0:1] == "v" {
		phonenumber := watoken.DecodeGetId(PublicKey, roomId)
		if phonenumber != "" {
			var infologin LoginInfo
			infologin.Uuid = roomId
			infologin.Login = roomId
			infologin.Phone = phonenumber
			SendStructTo(roomId, infologin)
		}
	}
}
