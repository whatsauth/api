package chatroot

import (
	"log"

	"github.com/gofiber/websocket/v2"
)

type Client struct{} // Add more data to this type if needed

var Clients = make(map[*websocket.Conn]Client) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
var Register = make(chan *websocket.Conn)
var Broadcast = make(chan string)
var Unregister = make(chan *websocket.Conn)

func RunHub() {
	for {
		select {
		case connection := <-Register:
			Clients[connection] = Client{}
			log.Println("connection registered")

		case message := <-Broadcast:
			log.Println("message received:", message)

			// Send the message to all clients
			for connection := range Clients {
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Println("write error:", err)

					Unregister <- connection
					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
				}
			}

		case connection := <-Unregister:
			// Remove the client from the hub
			delete(Clients, connection)

			log.Println("connection unregistered")
		}
	}
}

func RunSocket(c *websocket.Conn) {
	// When the function returns, unregister the client and close the connection
	defer func() {
		Unregister <- c
		c.Close()
	}()

	// Register the client
	Register <- c

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error:", err)
			}

			return // Calls the deferred function, i.e. closes the connection on error
		}

		if messageType == websocket.TextMessage {
			// Broadcast the received message
			Broadcast <- string(message)
		} else {
			log.Println("websocket message received of type", messageType)
		}
	}
}
