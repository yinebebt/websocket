//go:build client

package main

import (
	"golang.org/x/net/websocket"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	url := "ws://localhost:3000/channel"
	ws, err := websocket.Dial(url, "", "*")
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// sending a message to the server
	message := "Hello Server!"
	if _, err := ws.Write([]byte(message)); err != nil {
		log.Fatal(err)
	}
	log.Printf("sent:%s\n", message)

	// reading messages from the server
	go func() {
		for {
			var msg = make([]byte, 512)
			n, err := ws.Read(msg)
			if err != nil {
				log.Printf("read error:%v\n", err)
				return
			}
			log.Printf("received:%s\n", msg[:n])
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt
	log.Println("closing connection...")
}
