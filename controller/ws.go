package controller

import (
	"api/config"
	"api/helper/ws"

	"github.com/gofiber/websocket/v2"
)

func WsWhatsAuth(c *websocket.Conn) {
	ws.RunSocket(c, config.PublicKey)
}
