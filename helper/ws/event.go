package ws

import (
	"time"

	"api/helper/watoken"
)

func MagicLinkEvent(roomId string, PublicKey, PrivateKey string) {
	if roomId[0:1] == "v" {
		payload, err := watoken.Decode(PublicKey, roomId)
		if err == nil {
			newlogin, _ := watoken.EncodeforHours(payload.Id, payload.Alias, PrivateKey, 18)
			var infologin LoginInfo
			infologin.Uuid = roomId
			infologin.Login = newlogin
			infologin.Phone = payload.Id
			infologin.Alias = payload.Alias
			time.Sleep(1 * time.Millisecond)
			SendStructTo(roomId, infologin)

		}
	}
}
