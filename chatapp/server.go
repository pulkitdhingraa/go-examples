package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn	// connection for each client
	send chan []byte	//server pushes msg to send channel which deliver msg to client side through web socket connection
}

type Server struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewServer() *Server {
	return &Server{
		clients: make(map[*Client]bool),
		broadcast: make(chan []byte),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (s *Server) Run() {
	for {
		select {
		case client := <-s.register:
			s.clients[client] = true
			log.Println("New Client connected")
		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
				log.Println("Client disconnected")
			}
		case message := <-s.broadcast:
			for client := range s.clients {
				client.send <- message
			}
		}
	}
}

func (s *Server) HandleClient(conn *websocket.Conn) {
	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	s.register <- client

	defer func(){
		s.unregister <- client
		conn.Close()
	}()

	go func() {
		for message := range client.send {
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		formattedMsg := fmt.Sprintf("Client %p: %s", client, message)
		
		s.broadcast <- []byte(formattedMsg)	  
	}
}

func main() {
	port := flag.String("port", "8080", "port to start server on")
	flag.Parse()

	server := NewServer()
	go server.Run()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<- interrupt
		log.Println("Shutting down server...")
		os.Exit(0)
	}()

	upgrader := websocket.Upgrader{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wsConn, err := upgrader.Upgrade(w,r,nil)
		if err != nil {
			log.Println("Error upgrading Websocket:", err)
			return
		}

		go server.HandleClient(wsConn)
	})

	log.Printf("Server started on :%s", *port)
	if err := http.ListenAndServe("localhost:"+*port, nil); err != nil {
		log.Fatal("Listen and server error:", err)
	}
}