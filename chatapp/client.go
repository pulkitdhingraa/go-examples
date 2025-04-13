package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {
	serverAddr := flag.String("port", "localhost:8080", "server address")
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{
		Scheme: "ws",
		Host: *serverAddr,
		Path: "/",
	}
	log.Printf("Connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error connecting to websocket server:", err)
	}
	defer conn.Close()
	log.Println("Connected to websocket server")

	done := make(chan struct{})

	go readMessages(conn, done)

	writeMessages(conn, done, interrupt)
}

func readMessages(conn *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading from server:", err)
			return
		}
		fmt.Printf("< %s\n", message)
	}
}

func writeMessages(conn *websocket.Conn, done chan struct{}, interrupt chan os.Signal) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Type messages and press Enter to send (ctrl+c to quit):")

	for {
		select {
		case <- done:
			return
		case <- interrupt:
			log.Println("Interrupt Received. Closing connection")
			return
		default:
			if scanner.Scan() {
				message := scanner.Text()
				err := conn.WriteMessage(websocket.TextMessage, []byte(message))
				if err != nil {
					log.Println("Error writing to server:", err)
					return
				}
			}
		}
	}
}