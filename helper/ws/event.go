package ws

import (
	"time"

	"github.com/whatsauth/watoken"
)

func MagicLinkEvent(roomId string, PublicKey, PrivateKey string) {
	if roomId[0:1] == "v" {
		phonenumber := watoken.DecodeGetId(PublicKey, roomId)
		if phonenumber != "" {
			newlogin, _ := watoken.Encode(phonenumber, PrivateKey)
			var infologin LoginInfo
			infologin.Uuid = roomId
			infologin.Login = newlogin
			infologin.Phone = phonenumber
			time.Sleep(1 * time.Millisecond)
			SendStructTo(roomId, infologin)

		}
	}
}
