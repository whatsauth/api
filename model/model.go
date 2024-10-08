package model

import "go.mau.fi/whatsmeow/types"

type Header struct {
	Token string `reqHeader:"token"`
}

type Pantun struct {
	Notif  string   `bson:"notif"`
	Pantun []string `bson:"pantun"`
}

type Response struct {
	Response string `json:"response"`
}

type WaList struct {
	PhoneNumbers []types.IsOnWhatsAppResponse `json:"phonenumbers"`
}
