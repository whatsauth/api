package wa

import (
	"go.mau.fi/whatsmeow"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DocumentMessage struct {
	To        string `json:"to"`
	Base64Doc string `json:"base64doc"`
	Filename  string `json:"filename,omitempty"`
	Caption   string `json:"caption,omitempty"`
	IsGroup   bool   `json:"isgroup,omitempty"`
}

type ImageMessage struct {
	To          string `json:"to"`
	Base64Image string `json:"base64image"`
	Caption     string `json:"caption,omitempty"`
	IsGroup     bool   `json:"isgroup,omitempty"`
}

type PhoneList struct {
	PhoneNumbers []string `json:"phonenumbers"`
}

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

type Clients struct {
	List []*WaClient
}

type WaClient struct {
	PhoneNumber    string
	WAClient       *whatsmeow.Client
	eventHandlerID uint32
	Mongoconn      *mongo.Database
}

type WebHook struct {
	URL    string `bson:"url" json:"url"`
	Secret string `bson:"secret" json:"secret"`
}

type User struct {
	PhoneNumber   string  `bson:"phonenumber" json:"phonenumber"`
	DeviceID      uint16  `bson:"deviceid" json:"deviceid"`
	WebHook       WebHook `bson:"webhook" json:"webhook"`
	Mongostring   string  `bson:"mongostring" json:"mongostring"`
	Token         string  `bson:"token" json:"token"`
	ReadStatusOff bool    `bson:"readstatusoff" json:"readstatusoff"`
	SendTyping    bool    `bson:"sendtyping" json:"sendtyping"`
}

type LogSenderReceiver struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Sender   string             `bson:"sender"`
	Receiver string             `bson:"receiver"`
}

type LogSenderCounterUsage struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Sender  string             `bson:"sender"`
	Counter uint64             `bson:"counter,omitempty"`
	Users   []string           `bson:"users,omitempty"`
}
