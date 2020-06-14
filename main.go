package main

import (
	"fmt"
	"goIM/server"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	manager  = server.NewClientManager()
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	fmt.Println("Starting application...")
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	go manager.Start()
	http.HandleFunc("/ws", wsPage)
	http.ListenAndServe(":8080", nil)
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, error := upgrader.Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}

	client := server.NewClient(conn)

	manager.Register <- client

	go client.Read(manager)
	go client.Write()
}
