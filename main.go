package main

import (
	"fmt"
	"goIM/server"
	"goim/model"
	"net/http"

	"github.com/go-redis/redis"
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
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	// check redis connection
	fmt.Println("connect redis")
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	fmt.Println("Starting IM server")
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	go manager.Start()
	fmt.Println("Starting Websocket server")
	http.HandleFunc("/ws", wsPage)
	http.ListenAndServe(":8080", nil)
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, error := upgrader.Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}

	var msg *model.Message
	err := conn.ReadJSON(&msg)
	if err != nil {
		fmt.Println("login in error")
		return
	}

	client := server.NewClient(conn, msg.Username)

	manager.Register <- client

	go client.Read(manager)
	go client.Write()
}
