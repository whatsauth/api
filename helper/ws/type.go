package ws

import "github.com/gofiber/websocket/v2"

type Client struct {
	Id   string
	Conn *websocket.Conn
} // Register Conn socket with ID

type Message struct {
	Id      string
	Message string
} // To send message to Id

var Clients = make(map[string]*websocket.Conn) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
var Register = make(chan Client)               // Register channel for Client Struct
var SendMesssage = make(chan Message)
var Unregister = make(chan string)

type LoginInfo struct {
	Phone string `json:"phone,omitempty" bson:"phone,omitempty"`
	Login string `json:"login,omitempty" bson:"login,omitempty"`
	Uuid  string `json:"uuid,omitempty" bson:"uuid,omitempty"`
}

type WhatsauthRequest struct {
	Uuid        string `json:"uuid,omitempty" bson:"uuid,omitempty"`
	Phonenumber string `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
	Delay       uint32 `json:"delay,omitempty" bson:"delay,omitempty"`
}
