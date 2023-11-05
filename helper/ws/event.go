package ws

import (
	"log"
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
			n := 1
			for Clients[roomId] == nil && n > 10 {
				success := SendStructTo(roomId, infologin)
				time.Sleep(500 * time.Millisecond)
				log.Println(success)
				n++
			}
		}
	}
}
