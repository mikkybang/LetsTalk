package model

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/metaclips/FinalYearProject/values"
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
	c.WS.SetPongHandler(
		func(string) error {
			c.WS.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		},
	)

	for {
		_, msg, err := c.WS.ReadMessage()
		fmt.Println(string(msg))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v\n", err)
			}
			break
		}

		var data map[string]interface{}
		err = json.Unmarshal(msg, &data)
		if err != nil {
			// todo: handle marshal error better.
			return
		}
		msgType, ok := data["msgType"].(string)
		if !ok {
			log.Println("User did not send a valid message type", user, data)
			return
		}

		switch msgType {
		// Usin .(type) casting instead to convert data to required struct
		// doesnt seem to work as json UnMarshal converts data to a map
		// todo: any better way to cast
		// todo: add support to remove message.
		// todo: treat errors better.
		// todo: add support to remove messages.
		// todo: users should choose if to join chat.
		case "NewRoomCreated":
			var convertedType NewRoomRequest
			if err := json.Unmarshal(msg, &convertedType); err != nil {
				log.Println("Could not convert to required New Room Request struct")
				return
			}

			roomID, err := convertedType.CreateNewRoom()
			if err != nil {
				log.Println("Unable to create a new room for user:", convertedType.Email, "err:", err.Error())
				return
			}

			// Broadcast a joined message.
			userJoinedMessage := Joined{
				RoomID:      roomID,
				Email:       convertedType.Email,
				Joined:      true,
				RoomName:    convertedType.RoomName,
				MessageType: "UserJoinedRoom",
			}
			jsonByte, err := json.Marshal(userJoinedMessage)
			if err != nil {
				log.Println("Could not marshal to jsonByte while creating room", err.Error())
				return
			}
			m := WSMessage{jsonByte, convertedType.Email}
			HubConstruct.Broadcast <- m

		case "RequestUsersToJoinRoom":
			users, ok := data["users"].([]interface{})
			if ok {
				var convertedType JoinRequest
				if err := json.Unmarshal(msg, &convertedType); err != nil {
					log.Println("Could not convert to required Joined Request struct")
					return
				}

				for i := range users {
					user, ok := users[i].(string)
					if ok {
						roomRegisteredUser, err := convertedType.RequestUserToJoinRoom(user)
						if err != nil {
							log.Println("Error while requesting to room", err)
							continue
						}

						mapContent := map[string]interface{}{
							"requesterID":   convertedType.RequestingUserID,
							"requesterName": convertedType.RequestingUserName,
							"userRequested": user,
							"roomID":        convertedType.RoomID,
							"roomName":      convertedType.RoomName,
							"msgType":       "RequestUsersToJoinRoom",
						}

						jsonContent, err := json.Marshal(mapContent)
						if err != nil {
							log.Println("could not marshal to RequestUsersToJoinRoom, err:", err)
							return
						}

						// Send back RequestUsersToJoinRoom signal to everyone
						for _, roomRegisteredUser := range roomRegisteredUser {
							m := WSMessage{jsonContent, roomRegisteredUser}
							HubConstruct.Broadcast <- m
						}
						m := WSMessage{jsonContent, user}
						HubConstruct.Broadcast <- m
					}
				}

			} else {
				log.Println("could not convert users details to a []string")
			}
		case "UserJoinedRoom":
			var convertedType Joined
			if err := json.Unmarshal(msg, &convertedType); err != nil {
				log.Println("Could not convert to required Join Room Request struct")
				return
			}

			if convertedType.Email != user {
				return
			}
			users, err := convertedType.JoinRoom()
			if err != nil {
				return
			}

			for _, user := range users {
				m := WSMessage{msg, user}
				HubConstruct.Broadcast <- m
			}

		case "RequestAllMessages":
			roomID, ok := data["roomID"].(string)
			if ok {
				messages, roomName, err := GetAllMessageInRoom(roomID)
				if err != nil {
					log.Println("could not get all messages in room, err:", err)
					return
				}
				mapContent := map[string]interface{}{
					"messages": messages,
					"msgType":  "RequestAllMessages",
					"roomName": roomName,
				}

				jsonContent, err := json.Marshal(mapContent)
				if err != nil {
					log.Println("could not marshal images, err:", err)
					return
				}

				m := WSMessage{jsonContent, user}
				HubConstruct.Broadcast <- m
			}

		case "NewMessage":
			var convertedType Message
			if err := json.Unmarshal(msg, &convertedType); err != nil {
				log.Println("Could not convert to required New Message struct")
				return
			}

			if user != convertedType.User {
				return
			}
			convertedType.Time = time.Now().Format(values.TimeLayout)

			// Send message to all users.
			// Message is sent back to you as confirmation
			// it is delivered and saved to DB.
			registeredUsers, err := convertedType.SaveMessageContent()
			if err != nil {
				log.Println("Error saving msg to db", err, user)
				return
			}
			for _, registeredUser := range registeredUsers {
				m := WSMessage{msg, registeredUser}
				HubConstruct.Broadcast <- m
			}
		default:
			log.Println("Could not convert required type", msgType)
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
