package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options
var connections = make(map[*websocket.Conn]bool)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":3000", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// store the connection for later use 
	connections[c] = true

	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			onConnectionError(c, err)
			return
		}

		log.Printf("%s", message)

		for cc := range connections {
			if c != cc {
				err = cc.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					onConnectionError(c, err)
					return
				}
			}
		}
	}
}

func onConnectionError(c *websocket.Conn, err error) {
	log.Printf("error: %v", err)
	delete(connections, c)
}
