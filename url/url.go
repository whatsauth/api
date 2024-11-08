package url

import (
	"api/controller"
	"api/helper/chatroot"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Web(page *fiber.App) {
	//page.Get("/api/start/device/+", controller.StartDevice)
	page.Get("/api/device/+", controller.Device)
	page.Post("/api/numbers/isonwa", controller.CheckNumbersInWhatsApp)

	page.Post("/api/signup", controller.SignUp)
	//page.Post("/api/send/message/text", controller.SendTextMessage)
	//violation error v1 dimatikan dahulu karena terlalu banyak
	page.Post("/api/v2/send/message/text", controller.SendTextMessageV2)                  // kirim text message v2
	page.Post("/api/v3/official/send/message/text", controller.SendTextMessageV3FromUser) // kirim text message v3 dengan token adalah user official
	page.Post("/api/send/message/image", controller.SendImageMessage)
	page.Post("/api/send/message/document", controller.SendDocumentMessage)

	page.Post("/api/whatsauth/request", controller.PostWhatsAuthRequest)

	//websocket
	page.Get("/ws/whatsauth/public", websocket.New(controller.WsWhatsAuth)) //qr auth
	page.Get("/ws/chatgpl/public", websocket.New(chatroot.RunSocket))       //chatGPL public
	//page.Get("/ws/webrtc/public", websocket.New(wrtc.RunWebRTCSocket))      // New route for WebRTC signaling
}
