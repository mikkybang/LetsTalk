package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"github.com/metaclips/LetsTalk/values"
)

type messageBytes []byte

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

// WritePump pumps messages from the hub to the websocket connection.
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

// write writes a message with the given message type and payload.
func (c *Connection) write(mt int, payload []byte) error {
	if err := c.WS.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}

	return c.WS.WriteMessage(mt, payload)
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
			return c.WS.SetReadDeadline(time.Now().Add(pongWait))
		})

	for {
		var err error
		var msg messageBytes
		_, msg, err = c.WS.ReadMessage()
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
			log.Println("user did not send a valid message type", user, data)
			return
		}

		switch msgType {
		// todo: add support to remove message.
		// todo: add support to remove messages.
		// todo: users should choose if to join chat.
		case "NewRoomCreated":
			msg.handleCreateNewRoom()

		case "RequestUsersToJoinRoom":
			msg.handleRequestUserToJoinRoom()

		case "UserJoinedRoom":
			msg.handleUserAcceptRoomRequest(user)

		case "RequestAllMessages":
			roomID, ok := data["roomID"].(string)
			if ok {
				handleRequestAllMessages(roomID, user)
			}

		case "NewMessage":
			msg.handleNewMessage(user)
		default:
			log.Println("Could not convert required type", msgType)
		}
	}
}

func (msg messageBytes) handleCreateNewRoom() {
	var newRoom NewRoomRequest
	if err := json.Unmarshal(msg, &newRoom); err != nil {
		log.Println("Could not convert to required New Room Request struct")
		return
	}

	roomID, err := newRoom.CreateNewRoom()
	if err != nil {
		log.Println("Unable to create a new room for user:", newRoom.Email, "err:", err.Error())
		return
	}

	// Broadcast a joined message.
	userJoinedMessage := Joined{
		RoomID:      roomID,
		Email:       newRoom.Email,
		RoomName:    newRoom.RoomName,
		MessageType: "UserJoinedRoom",
	}

	jsonByte, err := json.Marshal(userJoinedMessage)
	if err != nil {
		log.Println("Could not marshal to jsonByte while creating room", err.Error())
		return
	}

	m := WSMessage{jsonByte, newRoom.Email}
	HubConstruct.Broadcast <- m
}

func (msg messageBytes) handleRequestUserToJoinRoom() {
	var request JoinRequest
	if err := json.Unmarshal(msg, &request); err != nil {
		log.Println("Could not convert to required Joined Request struct")
		return
	}

	for _, user := range request.Users {
		roomRegisteredUser, err := request.RequestUserToJoinRoom(user)
		if err != nil {
			log.Println("Error while requesting to room", err)
			continue
		}

		data := map[string]interface{}{
			"requesterID":   request.RequestingUserID,
			"requesterName": request.RequestingUserName,
			"userRequested": user,
			"roomID":        request.RoomID,
			"roomName":      request.RoomName,
			"msgType":       "RequestUsersToJoinRoom",
		}

		jsonContent, err := json.Marshal(data)
		if err != nil {
			log.Println("could not marshal to RequestUsersToJoinRoom, err:", err)
			continue
		}

		// Send back RequestUsersToJoinRoom signal to everyone registered in room.
		for _, roomRegisteredUser := range roomRegisteredUser {
			m := WSMessage{jsonContent, roomRegisteredUser}
			HubConstruct.Broadcast <- m
		}

		m := WSMessage{jsonContent, user}
		HubConstruct.Broadcast <- m
	}
}

func (msg messageBytes) handleUserAcceptRoomRequest(joiner string) {
	var roomRequest Joined
	if err := json.Unmarshal(msg, &roomRequest); err != nil {
		log.Println("Could not convert to required Join Room Request struct")
		return
	}

	if roomRequest.Email != joiner {
		return
	}

	users, err := roomRequest.AcceptRoomRequest()
	if err != nil {
		log.Println("could not join room", err)
		return
	}

	for _, user := range users {
		m := WSMessage{msg, user}
		HubConstruct.Broadcast <- m
	}
}

func handleRequestAllMessages(roomID, requester string) {
	messages, roomName, err := GetAllMessageInRoom(roomID)
	if err != nil {
		log.Println("could not get all messages in room, err:", err)
		return
	}

	data := map[string]interface{}{
		"messages": messages,
		"msgType":  "RequestAllMessages",
		"roomName": roomName,
		"roomID":   roomID,
	}

	jsonContent, err := json.Marshal(data)
	if err != nil {
		log.Println("could not marshal images, err:", err)
		return
	}

	m := WSMessage{jsonContent, requester}
	HubConstruct.Broadcast <- m
}

func (msg messageBytes) handleNewMessage(requester string) {
	var newMessage Message
	if err := json.Unmarshal(msg, &newMessage); err != nil {
		log.Println("Could not convert to required New Message struct")
		return
	}

	if requester != newMessage.UserID {
		return
	}
	newMessage.Time = time.Now().Format(values.TimeLayout)
	// Message is sent back to all users including sender.
	registeredUsers, err := newMessage.SaveMessageContent()
	if err != nil {
		log.Println("Error saving msg to db", err, requester)
		return
	}

	jsonContent, err := json.Marshal(newMessage)
	if err != nil {
		log.Println("Error converted message to json content", err)
		return
	}

	for _, registeredUser := range registeredUsers {
		m := WSMessage{jsonContent, registeredUser}
		HubConstruct.Broadcast <- m
	}
}
