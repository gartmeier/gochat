package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/websocket"
)

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://0.0.0.0:3000/", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt

		// Close the connection when an interrupt is received
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

	// Goroutine to continuously read messages from the server
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("error:", err)
				return
			}
			fmt.Printf("\033[2K\r< %s\n> ", message)
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		m, _ := reader.ReadString('\n')
		m = strings.TrimSuffix(m, "\n")

		err = c.WriteMessage(websocket.TextMessage, []byte(m))
		if err != nil {
			log.Fatal(err)
		}
	}
}

