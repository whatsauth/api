package ws

import (
	"log"
	"time"

	"github.com/whatsauth/watoken"
)

func MagicLinkEvent(roomId string, PublicKey, PrivateKey string) {
	log.Print("Masuk Ke Magic Link Event, status socket : ")
	if roomId[0:1] == "v" {
		phonenumber := watoken.DecodeGetId(PublicKey, roomId)
		if phonenumber != "" {
			newlogin, _ := watoken.Encode(phonenumber, PrivateKey)
			var infologin LoginInfo
			infologin.Uuid = roomId
			infologin.Login = newlogin
			infologin.Phone = phonenumber
			time.Sleep(1 * time.Second)
			success := SendStructTo(roomId, infologin)
			log.Print(success)
			for !success {
				success = SendStructTo(roomId, infologin)
				time.Sleep(1 * time.Second)
				log.Print(success)
			}
		}
	}
}
