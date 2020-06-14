package server

import (
	"goim/model"

	"github.com/gorilla/websocket"
)

// Client model
type Client struct {
	UserName string
	socket   *websocket.Conn
	Send     chan model.Message
}

func NewClient(conn *websocket.Conn, name string) (c *Client) {
	c = &Client{
		UserName: name,
		socket:   conn,
		Send:     make(chan model.Message),
	}
	return
}

func (c *Client) Read(clientManager *ClientManager) {
	defer func() {
		clientManager.Unregister <- c
		c.socket.Close()
	}()

	for {
		var msg model.Message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			clientManager.Unregister <- c
			c.socket.Close()
			break
		}
		//jsonMessage, _ := json.Marshal(&model.Message{Email: c.Message.Email, Content: string(message)})
		clientManager.Broadcast <- msg
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

			c.socket.WriteJSON(message)
		}
	}
}
