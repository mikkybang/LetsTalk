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
			if HubConstruct.Users[roomRegisteredUser] != nil {
				m := WSMessage{jsonContent, roomRegisteredUser}
				HubConstruct.Broadcast <- m
			}
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
		if HubConstruct.Users[user] != nil {
			m := WSMessage{msg, user}
			HubConstruct.Broadcast <- m
		}
	}
}

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
	fmt.Println(string(msg))
	if err := json.Unmarshal(msg, &newMessage); err != nil {
		log.Println("Could not convert to required New Message struct", err)
		return
	}

	if requester != newMessage.UserID {
		return
	}

	newMessage.Time = time.Now().Format(values.TimeLayout)
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

	if data, err := json.Marshal(request); err == nil {
		m := WSMessage{data, email}
		HubConstruct.Broadcast <- m
	}
}

func broadcastOnlineStatusToAllUserRoom(userEmail string, online bool) {
	// Update all users associates if online or not.
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
			// from HubRun, we should it in a goroutine so as to make broadcast
			HubConstruct.Broadcast <- m
		}
	}
}
