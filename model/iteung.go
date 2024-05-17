package model

import (
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type IteungV1Message struct {
	Phone_number string  `json:"phone_number"`
	Group_name   string  `json:"group_name"`
	Alias_name   string  `json:"alias_name"`
	Messages     string  `json:"messages"`
	Is_group     string  `json:"is_group"`
	Filename     string  `json:"filename"`
	Filedata     string  `json:"filedata"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Api_key      string  `json:"api_key"`
}

type IteungMessage struct {
	Phone_number       string  `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	Reply_phone_number string  `json:"reply_phone_number,omitempty" bson:"reply_phone_number,omitempty"`
	Chat_number        string  `json:"chat_number,omitempty" bson:"chat_number,omitempty"`
	Chat_server        string  `json:"chat_server,omitempty" bson:"chat_server,omitempty"`
	Group_name         string  `json:"group_name,omitempty" bson:"group_name,omitempty"`
	Group_id           string  `json:"group_id,omitempty" bson:"group_id,omitempty"`
	Group              string  `json:"group,omitempty" bson:"group,omitempty"`
	Alias_name         string  `json:"alias_name,omitempty" bson:"alias_name,omitempty"`
	Message            string  `json:"messages,omitempty" bson:"messages,omitempty"`
	From_link          bool    `json:"from_link,omitempty" bson:"from_link,omitempty"`
	From_link_delay    uint32  `json:"from_link_delay,omitempty" bson:"from_link_delay,omitempty"`
	Is_group           bool    `json:"is_group,omitempty" bson:"is_group,omitempty"`
	Filename           string  `json:"filename,omitempty" bson:"filename,omitempty"`
	Filedata           string  `json:"filedata,omitempty" bson:"filedata,omitempty"`
	Latitude           float64 `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude          float64 `json:"longitude,omitempty" bson:"longitude,omitempty"`
	LiveLoc            bool    `json:"liveloc,omitempty" bson:"liveloc,omitempty"`
}

type Module struct {
	Name    string   `json:"name,omitempty" bson:"name,omitempty"`
	Keyword []string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type Typo struct {
	From string `json:"from,omitempty" bson:"from,omitempty"`
	To   string `json:"to,omitempty" bson:"to,omitempty"`
}

type IteungWhatsMeowConfig struct {
	Info     *types.MessageInfo
	Message  *waProto.Message
	Waclient *whatsmeow.Client
}

type IteungDBConfig struct {
	MongoConn        *mongo.Database
	TypoCollection   string
	ModuleCollection string
}

type GowaNotif struct {
	User     string `json:"user,omitempty" bson:"user,omitempty"`
	Server   string `json:"server,omitempty" bson:"server,omitempty"`
	Messages string `json:"messages,omitempty" bson:"messages,omitempty"`
}
