package main

import (
	gameHTTP "github.com/DavidLoftus/highsociety/http"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error while upgrading to WebSocket: ", err)
		return
	}

	p := gameHTTP.NewPlayer(conn)
	go func() {
		if err := p.Handle(); err != nil {

		}
	}()
}

func main() {
	handler := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	handler.Handle("/", fs)

	handler.HandleFunc("/ws/", HandleWebsocket)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}
