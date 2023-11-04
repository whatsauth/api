package wa

import "go.mau.fi/whatsmeow"

type QRStatus struct {
	PhoneNumber string `json:"phonenumber"`
	Status      bool   `json:"status"`
	Code        string `json:"code"`
	Message     string `json:"message"`
}

type WaClient struct {
	PhoneNumber    string
	WAClient       *whatsmeow.Client
	eventHandlerID uint32
}
