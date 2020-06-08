package server

import (
	"encoding/json"
	"goim/model"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// Client model
type Client struct {
	id     string
	socket *websocket.Conn
	Send   chan []byte
}

func NewClient(conn *websocket.Conn) (c *Client) {
	uid, _ := uuid.NewV4()
	c = &Client{
		id:     uid.String(),
		socket: conn,
		Send:   make(chan []byte),
	}
	return
}

func (c *Client) Read(clientManager *ClientManager) {
	defer func() {
		clientManager.Unregister <- c
		c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			clientManager.Unregister <- c
			c.socket.Close()
			break
		}
		jsonMessage, _ := json.Marshal(&model.Message{Sender: c.id, Content: string(message)})
		clientManager.Broadcast <- jsonMessage
	}
}

func (c *Client) Write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
