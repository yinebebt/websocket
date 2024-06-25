/*websocket demo via Golang websocket*/
package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
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

func (s *Server) HandleWS(ws *websocket.Conn) {
	fmt.Println("New WebSocket Connection:", ws.RemoteAddr().String())
	s.mu.Lock()
	s.conns[ws] = true
	s.mu.Unlock()

	s.readLoop(ws)

	// Remove connection when done
	s.mu.Lock()
	delete(s.conns, ws)
	s.mu.Unlock()
	ws.Close()
	fmt.Println("WebSocket Connection closed:", ws.RemoteAddr().String())
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading from websocket:", err)
			break
		}
		msg := buf[:n]
		s.broadcast(msg)
		fmt.Println("Message received:", string(msg))
		if _, err := ws.Write([]byte("Thank you for the message!")); err != nil {
			fmt.Println("Error writing to websocket:", err)
			break
		}
	}
}

func (s *Server) broadcast(msg []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(msg); err != nil {
				fmt.Println("Error writing to websocket:", err)
			}
		}(ws)
	}
}

func handshake(config *websocket.Config, req *http.Request) error {
	// todo: Custom handshake logic (e.g., not checking Origin header)
	if config.Origin.String() != "" {
		req.Header.Set("Origin", config.Origin.String())
	}
	return nil
}

func main() {
	server := NewServer()
	ser := websocket.Server{
		Handler:   websocket.Handler(server.HandleWS),
		Handshake: handshake,
	}
	http.Handle("/ws", ser)
	// websocket.Handler handle handshake implicitly,it checks if the origin is valid url.
	// http.Handle("/ws", websocket.Handler(server.HandleWS))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
