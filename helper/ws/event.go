package ws

import (
	"api/helper/wa"
	"fmt"
	"log"

	"github.com/whatsauth/watoken"
)

func EventReadSocket(roomId string, PublicKey string) {
	phonenumber := watoken.DecodeGetId(PublicKey, roomId)
	if phonenumber != "" {
		infologin := wa.QRStatus{
			PhoneNumber: phonenumber,
		}
		log.Println("Info Login EventReadSocket ", infologin)
		SendStructTo(roomId, infologin)
	} else {
		fmt.Println("EventReadSocket: phonenumber is empty")
	}

}
