package url

import (
	"api/controller"

	"github.com/gofiber/fiber/v2"
)

func Web(page *fiber.App) {
	page.Get("/", controller.Homepage)
	page.Get("/device/+", controller.Device)
	page.Post("/send/message", controller.SendMessage)
}
