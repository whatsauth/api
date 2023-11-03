package helper

import "go.mau.fi/whatsmeow"

var Client []WaClient

type QRStatus struct {
	PhoneNumber string `json:"phonenumber"`
	Status      bool   `json:"status"`
	QRCode      string `json:"qrcode"`
	Message     string `json:"message"`
}

type WaClient struct {
	PhoneNumber    string
	WAClient       *whatsmeow.Client
	eventHandlerID uint32
}
