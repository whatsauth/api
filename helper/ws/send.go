package ws

import (
	"encoding/json"
	"log"
)

func SendStructTo(ID string, strc interface{}) (res bool) {
	b, err := json.Marshal(strc)
	if err != nil {
		log.Printf("SendStructTo Error: %s", err)
		return
	}
	return SendMessageTo(ID, string(b))
}

func SendMessageTo(ID string, msg string) (res bool) {
	m := Message{
		Id:      ID,
		Message: msg,
	}
	if Clients[ID] == nil {
		log.Printf("Clients[ID] == nil , with m : %s", m)
		res = false
	} else {
		SendMesssage <- m
		res = true
	}
	return
}
