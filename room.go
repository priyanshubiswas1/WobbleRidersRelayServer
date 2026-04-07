package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	code  string
	host  *Client
	guest *Client
	mu    sync.RWMutex
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func newRoom(code string) *Room {
	return &Room{
		code: code,
	}
}

func (r *Room) addClient(c *Client) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.host == nil {
		r.host = c
		return "host"

	}
	if r.guest == nil {
		r.guest = c
		return "guest"

	}
	return "full"
}

func (r *Room) removeClient(c *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.host == c {
		r.host = nil
	} else if r.guest == c {
		r.guest = nil
	}
}

func (r *Room) relayToHost(msg []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.host != nil {
		r.host.send <- msg
	}
}

func (r *Room) broadcast(msg []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.host != nil {
		r.host.send <- msg
	}
	if r.guest != nil {
		r.guest.send <- msg
	}
}
