package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	rooms   = make(map[string]*Room)
	roomsMu sync.RWMutex
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func getOrCreateNewRoom(code string) *Room {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	if r, ok := rooms[code]; ok {
		return r
	}
	r := newRoom(code)
	rooms[code] = r
	return r
}

func sendJSON(c *Client, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Println("Marshal error:", err)
		return
	}
	c.send <- data
}

func writerPump(c *Client) {
	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}

func readerPump(c *Client, room *Room, role string) {
	defer func() {
		room.removeClient(c)
		close(c.send)
		log.Printf("[%s] %s disconnected", room.code, role)
	}()

	log.Printf("[%s] %s connected", room.code, role)

	sendJSON(c, InfoMessage{
		Type:    "info",
		Message: "joined as " + role,
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var base struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(raw, &base); err != nil {
			log.Println("Bad message:", err)
			continue
		}

		switch base.Type {
		case "input":
			if role == "guest" {
				// guest input goes to host only
				room.relayToHost(raw)
			}
		case "state":
			// host broadcasts game state to both players
			room.broadcast(raw)
		default:
			log.Printf("[%s] unknown message type: %s", room.code, base.Type)
		}
	}
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	roomCode := r.URL.Query().Get("room")
	if roomCode == "" {
		http.Error(w, "missing room parameter", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error: ", err)
		return
	}

	room := getOrCreateNewRoom(roomCode)
	client := &Client{
		conn: conn,
		send: make(chan []byte, 64),
	}
	role := room.addClient(client)

	if role == "full" {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"error","message":"room full"}`))
		conn.Close()
		return
	}

	go writerPump(client)
	readerPump(client, room, role)

	// for {
	// 	messageType, message, err := conn.ReadMessage()
	// 	if err != nil {
	// 		log.Println("Client disconnected:", conn.RemoteAddr())
	// 		break
	// 	}

	// 	log.Printf("Received: %s", message)

	// 	if err := conn.WriteMessage(messageType, message); err != nil {
	// 		log.Println("Write error:", err)
	// 		break
	// 	}
	// }
}

func main() {
	http.HandleFunc("/ws", handleConnection)
	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
