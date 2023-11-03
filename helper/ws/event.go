package ws

import (
	"api/helper"
	"fmt"
	"log"

	"github.com/whatsauth/watoken"
)

func EventReadSocket(roomId string, PublicKey string) {
	phonenumber := watoken.DecodeGetId(PublicKey, roomId)
	if phonenumber != "" {
		infologin := helper.QRStatus{
			PhoneNumber: phonenumber,
		}
		log.Println("Info Login EventReadSocket ", infologin)
		SendStructTo(roomId, infologin)
	} else {
		fmt.Println("EventReadSocket: phonenumber is empty")
	}

}
