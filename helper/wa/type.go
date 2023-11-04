package wa

import (
	"go.mau.fi/whatsmeow"
	"go.mongodb.org/mongo-driver/mongo"
)

type TextMessage struct {
	To       string `json:"to"`
	IsGroup  bool   `json:"isgroup,omitempty"`
	Messages string `json:"messages"`
}

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
	Mongoconn      *mongo.Database
}

type WebHook struct {
	URL    string `bson:"url"`
	Secret string `bson:"secret"`
}

type User struct {
	PhoneNumber string  `bson:"phonenumber"`
	WebHook     WebHook `bson:"webhook"`
	Token       string  `bson:"token"`
}
