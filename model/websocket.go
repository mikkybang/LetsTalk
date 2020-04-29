package model

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var HubConstruct = Hub{
	Broadcast:  make(chan WSMessage),
	Register:   make(chan Subscription),
	UnRegister: make(chan Subscription),
	Users:      make(map[string]map[*Connection]bool),
}

func (h *Hub) Run() {
	for {
		select {
		case s := <-h.Register:
			connections := h.Users[s.User]
			if connections == nil {
				connections = make(map[*Connection]bool)
				h.Users[s.User] = connections
			}
			h.Users[s.User][s.Conn] = true
		case s := <-h.UnRegister:
			connections := h.Users[s.User]
			if connections != nil {
				if _, ok := connections[s.Conn]; ok {
					delete(connections, s.Conn)
					close(s.Conn.Send)
					if len(connections) == 0 {
						delete(h.Users, s.User)
					}
				}
			}
		case m := <-h.Broadcast:
			connections := h.Users[m.User]
			for c := range connections {
				select {
				case c.Send <- m.Data:
				default:
					close(c.Send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.Users, m.User)
					}
				}
			}
		}
	}
}

// ReadPump pumps messages from the websocket connection to the hub.
func (s Subscription) ReadPump(user string) {
	c := s.Conn
	defer func() {
		HubConstruct.UnRegister <- s
		c.WS.Close()
	}()

	c.WS.SetReadLimit(maxMessageSize)
	c.WS.SetReadDeadline(time.Now().Add(pongWait))
	c.WS.SetPongHandler(func(string) error { c.WS.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, msg, err := c.WS.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v\n", err)
			}
			break
		}

		var data interface{}
		err = json.Unmarshal(msg, data)
		if err != nil {
			// todo: handle marshal error better.
			return
		}

		switch e := data.(type) {
		// todo: add support to remove message.
		// todo: treat errors better.
		// todo: add support to remove messages.
		// todo: users should choose if to join chat.
		case Joined:
			if e.Email != user {
				return
			}
			users, err := e.JoinOrExitRoom()
			if err != nil {
				return
			}
			// todo: fix this
			// Normally users requested should be asked to join or NOT..
			// Broadcast room exit/join to other users..
			for _, user := range users {
				m := WSMessage{msg, user}
				HubConstruct.Broadcast <- m
			}
			fmt.Println(e)
		case Message:
			if user != e.User {
				return
			}

			// We read content to multiple users on the current chat.
			registeredUsers, err := e.SaveMessageContent()
			if err != nil {
				log.Println("Error saving msg to db", err, user)
				return
			}
			for _, user := range registeredUsers {
				m := WSMessage{msg, user}
				HubConstruct.Broadcast <- m
			}
		default:
			log.Println("Could not convert required type", e)
		}
	}
}

// write writes a message with the given message type and payload.
func (c *Connection) write(mt int, payload []byte) error {
	c.WS.SetWriteDeadline(time.Now().Add(writeWait))
	return c.WS.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *Subscription) WritePump() {
	c := s.Conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.WS.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			// Still send message to other clients...
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
