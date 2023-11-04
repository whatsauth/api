package url

import (
	"api/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Web(page *fiber.App) {
	page.Get("/", controller.Homepage)
	page.Get("/device/+", controller.Device)
	page.Post("/send/message", controller.SendTextMessage)

	page.Get("/ws/whatsauth/public", websocket.New(controller.WsWhatsAuth))
}
