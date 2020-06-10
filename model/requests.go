package model

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/metaclips/LetsTalk/values"
)

type messageBytes []byte

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

	jsonContent, err := json.Marshal(userJoinedMessage)
	if err != nil {
		log.Println("Could not marshal to jsonByte while creating room", err.Error())
		return
	}

	if HubConstruct.Users[newRoom.Email] != nil {
		m := WSMessage{jsonContent, newRoom.Email}
		HubConstruct.Broadcast <- m
	}
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

		data := struct {
			RequesterID   string `json:"requesterID"`
			RequesterName string `json:"requesterName"`
			UserRequested string `json:"userRequested"`
			RoomID        string `json:"roomID"`
			RoomName      string `json:"roomName"`
			MsgType       string `json:"msgType"`
		}{
			request.RequestingUserID, request.RequestingUserName,
			user, request.RoomID, request.RoomName, "RequestUsersToJoinRoom",
		}

		jsonContent, err := json.Marshal(data)
		if err != nil {
			log.Println("could not marshal to RequestUsersToJoinRoom, err:", err)
			continue
		}

		// Send back RequestUsersToJoinRoom signal to everyone registered in room.
		for _, roomRegisteredUser := range roomRegisteredUser {
			if HubConstruct.Users[roomRegisteredUser] != nil {
				m := WSMessage{jsonContent, roomRegisteredUser}
				HubConstruct.Broadcast <- m
			}
		}

		if HubConstruct.Users[user] != nil {
			m := WSMessage{jsonContent, user}
			HubConstruct.Broadcast <- m
		}
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
		if HubConstruct.Users[user] != nil {
			m := WSMessage{msg, user}
			HubConstruct.Broadcast <- m
		}
	}
}

// handleNewMessage broadcasts users message to all online users and also saves to database.
func (msg messageBytes) handleNewMessage(requester string) {
	var newMessage Message
	if err := json.Unmarshal(msg, &newMessage); err != nil {
		log.Println("Could not convert to required New Message struct", err)
		return
	}

	// Do not send if registered WS user is not same as message sender.
	if requester != newMessage.UserID {
		return
	}

	newMessage.Time = time.Now().Format(values.TimeLayout)
	// Save message to database ensuring user is registered to room.
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

	// Message is sent back to all online users including sender.
	for _, registeredUser := range registeredUsers {
		if HubConstruct.Users[registeredUser] != nil {
			m := WSMessage{jsonContent, registeredUser}
			HubConstruct.Broadcast <- m
		}
	}
}

func (msg messageBytes) handleExitRoom(requester string) {
	// data := make(map[string]interface{})
	data := struct {
		Email  string `json:"email"`
		RoomID string `json:"roomID"`
	}{}

	if err := json.Unmarshal(msg, &data); err != nil {
		log.Println("Could not retrieve json on exit room request", err)
		return
	}

	if requester != data.Email {
		return
	}

	user := User{Email: data.Email}
	registeredUsers, err := user.ExitRoom(data.RoomID)
	if err != nil {
		log.Println("Error exiting room", err)
		return
	}

	// Broadcast to all online users of a room exit.
	for _, registeredUser := range registeredUsers {
		if HubConstruct.Users[registeredUser] != nil {
			m := WSMessage{msg, registeredUser}
			HubConstruct.Broadcast <- m
		}
	}

	if HubConstruct.Users[requester] != nil {
		m := WSMessage{msg, requester}
		HubConstruct.Broadcast <- m
	}
}

// handleSearchUser returns registered users that match searchText
func handleSearchUser(searchText, user string) {
	data := struct {
		UsersFound []string
		msgType    string
	}{
		GetUser(searchText, user),
		"getUsers",
	}

	jsonContent, err := json.Marshal(&data)
	if err != nil {
		log.Println("Error while converting search user result to json", err)
		return
	}

	if HubConstruct.Users[user] != nil {
		m := WSMessage{jsonContent, user}
		HubConstruct.Broadcast <- m
	}
}

// handleRequestAllMessages coallates all messages in a particular room
func handleRequestAllMessages(roomID, requester string) {
	room := Room{RoomID: roomID}
	if err := room.GetAllMessageInRoom(); err != nil {
		log.Println("could not get all messages in room, err:", err)
		return
	}

	onlineUsers := make(map[string]bool)
	for _, user := range room.RegisteredUsers {
		if name, ok := values.MapEmailToName[user]; ok {
			nameAndEmail := fmt.Sprintf("%s (%s)", name, user)
			onlineUsers[nameAndEmail] = HubConstruct.Users[user] != nil
		}
	}

	data := map[string]interface{}{
		"messages":    room.Messages,
		"msgType":     "RequestAllMessages",
		"roomName":    room.RoomName,
		"roomID":      roomID,
		"onlineUsers": onlineUsers,
	}

	jsonContent, err := json.Marshal(&data)
	if err != nil {
		log.Println("could not marshal images, err:", err)
		return
	}

	// TODO: There's a check to see if user is truly online,
	// before sending broadcast this is to reduce the number
	// of requests to HUB worker.
	// We can remove this per if we increase number of worker.
	if HubConstruct.Users[requester] != nil {
		m := WSMessage{jsonContent, requester}
		HubConstruct.Broadcast <- m
	}
}

// handleLoadUserContent loads all users contents on page load.
// All rooms joined and users requests are loaded through WS.å
func handleLoadUserContent(email string) {
	userInfo := User{Email: email}
	if err := userInfo.GetAllUserRooms(); err != nil {
		log.Println("Could not fetch users room", email)
		return
	}

	request := map[string]interface{}{
		"msgType":  "WebsocketOpen",
		"rooms":    userInfo.RoomsJoined,
		"requests": userInfo.JoinRequest,
	}

	if data, err := json.Marshal(request); err == nil && HubConstruct.Users[email] != nil {
		m := WSMessage{data, email}
		HubConstruct.Broadcast <- m
	}
}

// broadcastOnlineStatusToAllUserRoom broadcasts users availability
// status to all users joined rooms. Status are broadcasted timely.
func broadcastOnlineStatusToAllUserRoom(userEmail string, online bool) {
	user := User{Email: userEmail}
	associates, err := user.GetAllUsersAssociates()
	if err != nil {
		log.Println("could not get users associate", err)
		return
	}

	for _, assassociateEmail := range associates {
		if HubConstruct.Users[assassociateEmail] == nil {
			continue
		}

		nameAndEmail := fmt.Sprintf("%s (%s)", values.MapEmailToName[assassociateEmail], userEmail)
		msg := map[string]interface{}{
			"msgType":  "OnlineStatus",
			"username": nameAndEmail,
			"status":   online,
		}

		if data, err := json.Marshal(msg); err == nil {
			m := WSMessage{data, assassociateEmail}
			// Since we are calling broadcastOnlineStatusToAllUserRoom
			// from HubRun, we should call it in a goroutine so as
			// not to block the hub channel
			HubConstruct.Broadcast <- m
		}
	}
}
