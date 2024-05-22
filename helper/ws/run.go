package ws

import (
	"github.com/gofiber/websocket/v2"
)

func RunHub() { // Call this function on your main function before run fiber
	for {
		select {
		case connection := <-Register:
			Clients[connection.Id] = connection.Conn
		case message := <-SendMesssage:
			connection := Clients[message.Id]
			connection.WriteMessage(websocket.TextMessage, []byte(message.Message))
		case connection := <-Unregister:
			delete(Clients, connection)
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
			//log.Println("read error:", err)
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
			messageType, _, err := s.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					//log.Println("read error:", err)
				}
				return // Calls the deferred function, i.e. closes the connection on error
			}

			if messageType == websocket.TextMessage {
				// log the received message
				//log.Println(string(message))
			} else {
				//log.Println("websocket message received of type", messageType)
			}
		}
	} else {
		//log.Println("websocket message received of type", messageType)
	}
	return

}
