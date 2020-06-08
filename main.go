package main

import (
	"fmt"
	"goIM/server"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	manager = server.NewClientManager()
)

func main() {
	fmt.Println("Starting application...")
	manager := server.NewClientManager()
	go manager.Start()
	http.HandleFunc("/ws", wsPage)
	http.ListenAndServe(":12345", nil)
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}

	client := server.NewClient(conn)

	manager.Register <- client

	go client.Read(manager)
	go client.Write()
}
