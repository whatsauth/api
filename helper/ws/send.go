package ws

import (
	"encoding/json"
)

func SendStructTo(ID string, strc interface{}) (res bool) {
	b, err := json.Marshal(strc)
	if err != nil {
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
		res = false
	} else {
		SendMesssage <- m
		res = true
	}
	return
}
