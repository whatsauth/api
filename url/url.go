package url

import (
	"api/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Web(page *fiber.App) {
	//page.Get("/", controller.Homepage)
	page.Get("/api/device/+", controller.Device)
	page.Post("/api/send/message/text", controller.SendTextMessage)

	page.Post("/api/whatsauth/request", controller.PostWhatsAuthRequest)
	page.Get("/ws/whatsauth/public", websocket.New(controller.WsWhatsAuth))
}
