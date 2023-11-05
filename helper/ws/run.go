package ws

import (
	"log"

	"github.com/gofiber/websocket/v2"
)

func RunHub() { // Call this function on your main function before run fiber
	for {
		select {
		case connection := <-Register:
			Clients[connection.Id] = connection.Conn
			log.Println("connection registered : ")
			log.Println(connection)

		case message := <-SendMesssage:
			log.Println("message received:", message)
			connection := Clients[message.Id]
			err := connection.WriteMessage(websocket.TextMessage, []byte(message.Message))
			if err != nil {
				log.Println(err)
			}

		case connection := <-Unregister:
			// Remove the client from the hub
			delete(Clients, connection)

			log.Println("connection unregistered")
			log.Println(connection)
		}
	}
}

func RunSocket(c *websocket.Conn, PublicKey, PrivateKey string) (Id string) { // call this function after declare URL routes
	var s Client
	// When the function returns, unregister the client and close the connection
	defer func() {
		Unregister <- s.Id
		c.Close()
	}()
	messageType, message, err := c.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Println("read error:", err)
		}
		return // Calls the deferred function, i.e. closes the connection on error
	}
	Id = string(message)
	if messageType == websocket.TextMessage {
		// Get the received message
		// Register the client
		s = Client{
			Id:   Id,
			Conn: c,
		}
		Register <- s
		MagicLinkEvent(Id, PublicKey, PrivateKey)
		for {
			messageType, message, err := s.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}
				return // Calls the deferred function, i.e. closes the connection on error
			}

			if messageType == websocket.TextMessage {
				// log the received message
				log.Println(string(message))
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}
	} else {
		log.Println("websocket message received of type", messageType)
	}
	return

}
