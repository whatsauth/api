package helper

import "go.mau.fi/whatsmeow"

type QRStatus struct {
	Status  bool   `json:"status"`
	QRCode  string `json:"qrcode"`
	Message string `json:"message"`
}

type MyClient struct {
	WAClient       *whatsmeow.Client
	eventHandlerID uint32
}
