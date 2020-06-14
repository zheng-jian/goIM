package server

import (
	"goim/model"
)

//ClientManager struct
type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan model.Message
	Register   chan *Client
	Unregister chan *Client
}

func NewClientManager() (manager *ClientManager) {
	manager = &ClientManager{
		Broadcast:  make(chan model.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
	return
}

func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.Register:
			manager.Clients[conn] = true
			message := &model.Message{Message: "A new client has connected."}
			manager.send(*message, conn)

		case conn := <-manager.Unregister:
			if _, ok := manager.Clients[conn]; ok {
				close(conn.Send)
				delete(manager.Clients, conn)
				message := &model.Message{Message: "A client has disconnected."}
				manager.send(*message, conn)
			}
		case message := <-manager.Broadcast:
			for conn := range manager.Clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(manager.Clients, conn)
				}
			}
		}
	}
}

func (manager *ClientManager) send(message model.Message, ignore *Client) {
	for conn := range manager.Clients {
		//if conn != ignore {
		conn.Send <- message
		//}
	}
}
