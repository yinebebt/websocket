//go:build server

package main

import (
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type Server struct {
	conns map[*websocket.Conn]bool
	mu    sync.Mutex
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

var sessionCount int

func (s *Server) HandleWS(ws *websocket.Conn) {
	sessionCount++
	log.Printf("new session-%v started\n", sessionCount)
	s.mu.Lock()
	s.conns[ws] = true
	s.mu.Unlock()

	s.readLoop(ws)

	// Remove connection when done
	s.mu.Lock()
	delete(s.conns, ws)
	s.mu.Unlock()
	ws.Close()
	sessionCount--
	log.Println("ws session closed")
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("read error:%v", err)
			break
		}
		msg := buf[:n]
		if _, err := ws.Write(msg); err != nil {
			log.Printf("write error:%v\n", err)
			break
		}
	}
}

// broadcast sent back received message to all client.
func (s *Server) broadcast(msg []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(msg); err != nil {
				log.Printf("broadcast: write error:%v\n", err)
			}
		}(ws)
	}
}

func handshake(config *websocket.Config, req *http.Request) error {
	if origin := config.Origin; origin != nil {
		if origin.String() != req.Header.Get("Origin") {
			log.Printf("origin not allowed:%v\n", origin)
		}
		req.Header.Set("Origin", config.Origin.String())
	}
	return nil
}

func main() {
	server := NewServer()
	srv := websocket.Server{
		Handler:   websocket.Handler(server.HandleWS),
		Handshake: handshake,
		Config:    websocket.Config{},
	}
	http.Handle("/channel", srv)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("starting ws server at:%v/ws\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Printf("error while starting ws server:%v\n", err)
		return
	}
}
