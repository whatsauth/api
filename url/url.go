package url

import (
	"api/controller"
	"api/helper/chatroot"
	"api/helper/wrtc"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Web(page *fiber.App) {
	page.Get("/api/start/device/+", controller.StartDevice)
	page.Get("/api/device/+", controller.Device)

	page.Post("/api/signup", controller.SignUp)
	page.Post("/api/send/message/text", controller.SendTextMessage)
	page.Post("/api/send/message/image", controller.SendImageMessage)

	page.Post("/api/whatsauth/request", controller.PostWhatsAuthRequest)

	//websocket
	page.Get("/ws/whatsauth/public", websocket.New(controller.WsWhatsAuth)) //qr auth
	page.Get("/ws/chatgpl/public", websocket.New(chatroot.RunSocket))       //chatGPL public
	page.Get("/ws/webrtc/public", websocket.New(wrtc.RunWebRTCSocket))      // New route for WebRTC signaling
}
