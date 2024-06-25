package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	origin := "http://localhost/"
	url := "ws://localhost:3000/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Sending a message to the server
	message := "Hello Server!"
	if _, err := ws.Write([]byte(message)); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Sent:", message)

	// Reading messages from the server
	go func() {
		for {
			var msg = make([]byte, 512)
			n, err := ws.Read(msg)
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			fmt.Println("Received:", string(msg[:n]))
		}
	}()

	// Wait for an interrupt signal to gracefully shut down the client
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt
	fmt.Println("Interrupt received, closing connection...")
}
